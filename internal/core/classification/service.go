package classification

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/errx"
	"github.com/netbill/places-svc/internal/models"
	"github.com/netbill/restkit/pagi"
)

type Service struct {
	repo      repo
	media     media
	messenger messenger
	tx        transaction
}

type ServiceDeps struct {
	Repo      repo
	Media     media
	Messenger messenger
	Tx        transaction
}

func NewPlaceClassModule(deps ServiceDeps) *Service {
	return &Service{
		repo:      deps.Repo,
		media:     deps.Media,
		messenger: deps.Messenger,
		tx:        deps.Tx,
	}
}

type messenger interface {
	PublishPlaceClassCreated(ctx context.Context, class models.PlaceClass) error
	PublishPlaceClassUpdated(ctx context.Context, class models.PlaceClass) error
}

type media interface {
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

type CreateParams struct {
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	IconKey     *string    `json:"icon_key,omitempty"`
}

func (s *Service) Create(
	ctx context.Context,
	params CreateParams,
) (class models.PlaceClass, err error) {
	if params.ParentID != nil {
		class, err = s.repo.Get(ctx, *params.ParentID)
		if err != nil {
			return models.PlaceClass{}, err
		}

		if class.DeprecatedAt != nil {
			return models.PlaceClass{}, errx.ErrorPlaceClassIsDeprecated.Raise(
				fmt.Errorf("parent place class %s is deprecated", *params.ParentID),
			)
		}
	}

	if err = s.tx.Transaction(ctx, func(ctx context.Context) error {
		class, err = s.repo.Create(ctx, params)
		if err != nil {
			return err
		}

		return s.messenger.PublishPlaceClassCreated(ctx, class)
	}); err != nil {
		return models.PlaceClass{}, err
	}

	return class, nil
}

func (s *Service) Get(ctx context.Context, id uuid.UUID) (models.PlaceClass, error) {
	class, err := s.repo.Get(ctx, id)
	if err != nil {
		return models.PlaceClass{}, err
	}

	return class, nil
}

type FilterParams struct {
	Parent      *FilterParentParams
	BestMatch   *string
	Description *string
	Deprecated  *bool
}

type FilterParentParams struct {
	ID               uuid.UUID
	IncludedParents  bool
	IncludedChildren bool
}

func (s *Service) GetList(
	ctx context.Context,
	params FilterParams,
	limit, offset uint,
) (pagi.Page[[]models.PlaceClass], error) {
	classes, err := s.repo.GetList(ctx, params, limit, offset)
	if err != nil {
		return pagi.Page[[]models.PlaceClass]{}, err
	}

	return classes, nil
}

func (s *Service) GetByIDs(
	ctx context.Context,
	ids []uuid.UUID,
) ([]models.PlaceClass, error) {
	res, err := s.repo.GetListByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	return res, nil
}

type UpdateParams struct {
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	IconKey     *string    `json:"icon_key,omitempty"`
}

func (p UpdateParams) HasChanges(class models.PlaceClass) bool {
	return !ptrEqual(p.ParentID, class.ParentID) ||
		!ptrEqual(p.Name, &class.Name) ||
		!ptrEqual(p.Description, &class.Description) ||
		!ptrEqual(p.IconKey, class.IconKey)
}

func ptrEqual[T comparable](a, b *T) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

func (s *Service) Update(
	ctx context.Context,
	classID uuid.UUID,
	params UpdateParams,
) (class models.PlaceClass, err error) {
	class, err = s.repo.Get(ctx, classID)
	if err != nil {
		return models.PlaceClass{}, err
	}

	if params.ParentID != nil || *params.ParentID != uuid.Nil {
		exist, err := s.repo.CheckParentCycle(ctx, class.ID, *class.ParentID)
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
		if err := s.media.DeleteUploadPlaceClassIcon(ctx, class.ID, *class.IconKey); err != nil {
			return models.PlaceClass{}, err
		}
		params.IconKey = nil
	case params.IconKey != nil:
		iconKey, err := s.media.UpdatePlaceClassIcon(ctx, class.ID, *params.IconKey)
		if err != nil {
			return models.PlaceClass{}, err
		}
		params.IconKey = &iconKey
	}

	if err = s.tx.Transaction(ctx, func(ctx context.Context) error {
		class, err = s.repo.Update(ctx, classID, params)
		if err != nil {
			return err
		}

		return s.messenger.PublishPlaceClassUpdated(ctx, class)
	}); err != nil {
		return models.PlaceClass{}, err
	}

	return class, nil
}

func (s *Service) Deprecate(
	ctx context.Context,
	classID uuid.UUID,
) (models.PlaceClass, error) {
	return s.updateDeprecate(ctx, classID, true)
}

func (s *Service) Undeprecate(
	ctx context.Context,
	classID uuid.UUID,
) (models.PlaceClass, error) {
	return s.updateDeprecate(ctx, classID, false)
}

func (s *Service) updateDeprecate(
	ctx context.Context,
	classID uuid.UUID,
	value bool,
) (models.PlaceClass, error) {
	class, err := s.repo.Get(ctx, classID)
	if err != nil {
		return models.PlaceClass{}, err
	}

	if class.DeprecatedAt != nil && value || class.DeprecatedAt == nil && !value {
		return class, nil
	}

	if err = s.tx.Transaction(ctx, func(ctx context.Context) error {
		class, err = s.repo.Deprecated(ctx, classID, value)
		if err != nil {
			return err
		}

		return s.messenger.PublishPlaceClassUpdated(ctx, class)
	}); err != nil {
		return models.PlaceClass{}, err
	}

	return class, nil
}
