package pclass

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/models"
	"github.com/netbill/restkit/pagi"
)

type PlaceClassModule struct {
	repo      placeClassRepo
	media     placeClassBucket
	messenger placeClassMessenger
	tx        transaction
}

type transaction interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type PlaceClassModuleDeps struct {
	Repo      placeClassRepo
	Media     placeClassBucket
	Messenger placeClassMessenger
	tx        transaction
}

func NewPlaceClassModule(deps PlaceClassModuleDeps) *PlaceClassModule {
	return &PlaceClassModule{
		repo:      deps.Repo,
		media:     deps.Media,
		messenger: deps.Messenger,
		tx:        deps.tx,
	}
}

type placeClassRepo interface {
	CreatePlaceClass(ctx context.Context, class CreatePlaceClassParams) (models.PlaceClass, error)

	GetPlaceClass(
		ctx context.Context,
		id uuid.UUID,
	) (models.PlaceClass, error)
	GetPlaceClassesByIDs(
		ctx context.Context,
		ids []uuid.UUID,
	) ([]models.PlaceClass, error)
	GetPlaceClasses(
		ctx context.Context,
		params FilterParams,
		limit, offset uint,
	) (pagi.Page[[]models.PlaceClass], error)
	PlaceClassExists(
		ctx context.Context,
		classID uuid.UUID,
	) (bool, error)

	UpdatePlaceClass(
		ctx context.Context,
		classID uuid.UUID,
		params UpdateParams,
	) (models.PlaceClass, error)
	DeprecatedPlaceClass(
		ctx context.Context,
		classID uuid.UUID,
		value bool,
	) (models.PlaceClass, error)

	CheckParentCycle(ctx context.Context, classID, parentID uuid.UUID) (bool, error)
	CheckPlaceClassHasChildren(ctx context.Context, classID uuid.UUID) (bool, error)
}

type placeClassMessenger interface {
	PublishPlaceClassCreated(ctx context.Context, class models.PlaceClass) error
	PublishPlaceClassUpdated(ctx context.Context, class models.PlaceClass) error
}

type placeClassBucket interface {
	CreatePlaceClassUploadMediaLinks(
		ctx context.Context,
		classID uuid.UUID,
	) (models.UploadPlaceClassMediaLinks, error)

	UpdatePlaceClassIcon(
		ctx context.Context,
		classID uuid.UUID,
		oldKey *string,
		tempKey *string,
	) (newKey *string, err error)

	DeletePlaceClassIcon(
		ctx context.Context,
		classID uuid.UUID,
		key string,
	) error
}

type CreatePlaceClassParams struct {
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	IconKey     *string    `json:"icon_key,omitempty"`
}

func (m *PlaceClassModule) Create(
	ctx context.Context,
	params CreatePlaceClassParams,
) (class models.PlaceClass, err error) {
	if params.ParentID != nil {
		class, err = m.repo.GetPlaceClass(ctx, *params.ParentID)
		if err != nil {
			return models.PlaceClass{}, err
		}

		if class.DeprecatedAt != nil {
			return models.PlaceClass{}, errx.ErrorPlaceClassIsDeprecated.Raise(
				fmt.Errorf("parent place class %s is deprecated", *params.ParentID),
			)
		}
	}

	if err = m.tx.Transaction(ctx, func(ctx context.Context) error {
		class, err = m.repo.CreatePlaceClass(ctx, params)
		if err != nil {
			return err
		}

		return m.messenger.PublishPlaceClassCreated(ctx, class)
	}); err != nil {
		return models.PlaceClass{}, err
	}

	return class, nil
}

func (m *PlaceClassModule) Get(ctx context.Context, id uuid.UUID) (models.PlaceClass, error) {
	class, err := m.repo.GetPlaceClass(ctx, id)
	if err != nil {
		return models.PlaceClass{}, err
	}

	return class, nil
}

type FilterParams struct {
	Parent      *FilterPlaceClassParams
	BestMatch   *string
	Description *string
	Deprecated  *bool
}

type FilterPlaceClassParams struct {
	ID               uuid.UUID
	IncludedParents  bool
	IncludedChildren bool
}

func (m *PlaceClassModule) GetList(
	ctx context.Context,
	params FilterParams,
	limit, offset uint,
) (pagi.Page[[]models.PlaceClass], error) {
	classes, err := m.repo.GetPlaceClasses(ctx, params, limit, offset)
	if err != nil {
		return pagi.Page[[]models.PlaceClass]{}, err
	}

	return classes, nil
}

func (m *PlaceClassModule) GetByIDs(
	ctx context.Context,
	ids []uuid.UUID,
) ([]models.PlaceClass, error) {
	res, err := m.repo.GetPlaceClassesByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	return res, nil
}
