package core

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
	media     placeClassMedia
	messenger placeClassMessenger
	tx        transaction
}

type PlaceClassModuleDeps struct {
	Repo      placeClassRepo
	Media     placeClassMedia
	Messenger placeClassMessenger
	Tx        transaction
}

func NewPlaceClassModule(deps PlaceClassModuleDeps) *PlaceClassModule {
	return &PlaceClassModule{
		repo:      deps.Repo,
		media:     deps.Media,
		messenger: deps.Messenger,
		tx:        deps.Tx,
	}
}

type placeClassRepo interface {
	Create(ctx context.Context, params CreatePlaceClassParams) (models.PlaceClass, error)

	Get(
		ctx context.Context,
		id uuid.UUID,
	) (models.PlaceClass, error)
	GetListByIDs(
		ctx context.Context,
		ids []uuid.UUID,
	) ([]models.PlaceClass, error)
	GetList(
		ctx context.Context,
		params FilterParams,
		limit, offset uint,
	) (pagi.Page[[]models.PlaceClass], error)
	Exists(
		ctx context.Context,
		classID uuid.UUID,
	) (bool, error)

	Update(
		ctx context.Context,
		classID uuid.UUID,
		params UpdatePlaceClassParams,
	) (models.PlaceClass, error)
	Deprecated(
		ctx context.Context,
		classID uuid.UUID,
		value bool,
	) (models.PlaceClass, error)

	CheckParentCycle(ctx context.Context, classID, parentID uuid.UUID) (bool, error)
	CheckHasChildren(ctx context.Context, classID uuid.UUID) (bool, error)
}

type placeClassMessenger interface {
	PublishPlaceClassCreated(ctx context.Context, class models.PlaceClass) error
	PublishPlaceClassUpdated(ctx context.Context, class models.PlaceClass) error
}

type placeClassMedia interface {
	CreatePlaceClassIconUploadMediaLinks(
		ctx context.Context,
		classID uuid.UUID,
	) (models.UploadMediaLink, error)

	UpdatePlaceClassIcon(
		ctx context.Context,
		classID uuid.UUID,
		key string,
	) (string, error)

	DeleteUploadPlaceClassIcon(
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
		class, err = m.repo.Get(ctx, *params.ParentID)
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
		class, err = m.repo.Create(ctx, params)
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
	class, err := m.repo.Get(ctx, id)
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
	classes, err := m.repo.GetList(ctx, params, limit, offset)
	if err != nil {
		return pagi.Page[[]models.PlaceClass]{}, err
	}

	return classes, nil
}

func (m *PlaceClassModule) GetByIDs(
	ctx context.Context,
	ids []uuid.UUID,
) ([]models.PlaceClass, error) {
	res, err := m.repo.GetListByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	return res, nil
}

type UpdatePlaceClassParams struct {
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	IconKey     *string    `json:"icon_key,omitempty"`
}

func (P UpdatePlaceClassParams) HasChanges(class models.PlaceClass) bool {
	return !ptrEqual(P.ParentID, class.ParentID) ||
		!ptrEqual(P.Name, &class.Name) ||
		!ptrEqual(P.Description, &class.Description) ||
		!ptrEqual(P.IconKey, class.IconKey)
}

func ptrEqual[T comparable](a, b *T) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

func (m *PlaceClassModule) Update(
	ctx context.Context,
	classID uuid.UUID,
	params UpdatePlaceClassParams,
) (class models.PlaceClass, err error) {
	class, err = m.repo.Get(ctx, classID)
	if err != nil {
		return models.PlaceClass{}, err
	}

	if params.ParentID != nil || *params.ParentID != uuid.Nil {
		exist, err := m.repo.CheckParentCycle(ctx, class.ID, *class.ParentID)
		if err != nil {
			return models.PlaceClass{}, err
		}
		if exist {
			return models.PlaceClass{}, errx.ErrorPlaceClassParentCycle.Raise(
				fmt.Errorf("setting parent %s for class %s would create a cycle", *class.ParentID, class.ID),
			)
		}
	}

	switch {
	case params.IconKey != nil && *params.IconKey == "" && class.IconKey != nil:
		if err := m.media.DeleteUploadPlaceClassIcon(ctx, class.ID, *class.IconKey); err != nil {
			return models.PlaceClass{}, err
		}
		params.IconKey = nil
	case params.IconKey != nil:
		iconKey, err := m.media.UpdatePlaceClassIcon(ctx, class.ID, *params.IconKey)
		if err != nil {
			return models.PlaceClass{}, err
		}
		params.IconKey = &iconKey
	}

	if err = m.tx.Transaction(ctx, func(ctx context.Context) error {
		class, err = m.repo.Update(ctx, classID, params)
		if err != nil {
			return err
		}

		return m.messenger.PublishPlaceClassUpdated(ctx, class)
	}); err != nil {
		return models.PlaceClass{}, err
	}

	return class, nil
}

func (m *PlaceClassModule) Deprecate(
	ctx context.Context,
	classID uuid.UUID,
) (models.PlaceClass, error) {
	return m.updateDeprecate(ctx, classID, true)
}

func (m *PlaceClassModule) Undeprecate(
	ctx context.Context,
	classID uuid.UUID,
) (models.PlaceClass, error) {
	return m.updateDeprecate(ctx, classID, false)
}

func (m *PlaceClassModule) updateDeprecate(
	ctx context.Context,
	classID uuid.UUID,
	value bool,
) (models.PlaceClass, error) {
	class, err := m.repo.Get(ctx, classID)
	if err != nil {
		return models.PlaceClass{}, err
	}

	if class.DeprecatedAt != nil && value || class.DeprecatedAt == nil && !value {
		return class, nil
	}

	if err = m.tx.Transaction(ctx, func(ctx context.Context) error {
		class, err = m.repo.Deprecated(ctx, classID, value)
		if err != nil {
			return err
		}

		return m.messenger.PublishPlaceClassUpdated(ctx, class)
	}); err != nil {
		return models.PlaceClass{}, err
	}

	return class, nil
}
