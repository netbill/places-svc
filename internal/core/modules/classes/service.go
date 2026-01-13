package classes

import (
	"context"

	"github.com/google/uuid"
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

	CheckParentCycle(ctx context.Context, classID, parentID uuid.UUID) (bool, error)
	CheckClassHasChildren(ctx context.Context, classID uuid.UUID) (bool, error)
	CheckPlaceExistForClass(ctx context.Context, classID uuid.UUID) (bool, error)

	UpdatePlaceClassParent(ctx context.Context, classID uuid.UUID, parentID *uuid.UUID) (models.PlaceClass, error)
	DeletePlaceClass(ctx context.Context, classID uuid.UUID) error

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type messanger interface {
	PublishPlaceClassCreated(ctx context.Context, class models.PlaceClass) error
	PublishPlaceClassUpdated(ctx context.Context, class models.PlaceClass) error
	PublishPlaceClassParentUpdated(ctx context.Context, class models.PlaceClass) error
	PublishPlaceClassDeleted(ctx context.Context, classID uuid.UUID) error
}
