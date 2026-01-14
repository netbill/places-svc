package pclass

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/pagi"
	"github.com/netbill/places-svc/internal/core/models"
)

type Service struct {
	repo      repo
	messanger messanger
}

func New(repo repo, messanger messanger) Service {
	return Service{
		repo:      repo,
		messanger: messanger,
	}
}

type repo interface {
	CreatePlaceClass(ctx context.Context, class CreateParams) (models.PlaceClass, error)
	UpdatePlaceClass(ctx context.Context, classID uuid.UUID, params UpdateParams) (models.PlaceClass, error)

	GetPlaceClass(ctx context.Context, id uuid.UUID) (models.PlaceClass, error)
	GetPlaceClassByCode(ctx context.Context, code string) (models.PlaceClass, error)
	GetPlaceClasses(ctx context.Context, params FilterParams, limit, offset uint) (pagi.Page[[]models.PlaceClass], error)

	PlaceClassExists(ctx context.Context, classID uuid.UUID) (bool, error)
	PlaceClassExistsByCode(ctx context.Context, code string) (bool, error)

	CheckParentCycle(ctx context.Context, classID, parentID uuid.UUID) (bool, error)
	CheckPlaceClassHasChildren(ctx context.Context, classID uuid.UUID) (bool, error)
	CheckPlaceExistForClass(ctx context.Context, classID uuid.UUID) (bool, error)

	ReplacePlacesClassID(ctx context.Context, oldClassID, newClassID uuid.UUID) error

	DeletePlaceClass(ctx context.Context, classID uuid.UUID) error

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type messanger interface {
	PublishPlaceClassCreated(ctx context.Context, class models.PlaceClass) error
	PublishPlaceClassUpdated(ctx context.Context, class models.PlaceClass) error
	PublishPlaceClassParentUpdated(ctx context.Context, class models.PlaceClass) error
	PublishPlaceClassDeleted(ctx context.Context, classID uuid.UUID) error

	PublishPlacesClassReplaced(ctx context.Context, oldClassID, newClassID uuid.UUID) error
}
