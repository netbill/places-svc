package pclass

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/models"
	"github.com/netbill/restkit/pagi"
)

type Module struct {
	repo      repo
	bucket    bucket
	messanger messanger
	token     token
}

func New(repo repo, messanger messanger, bucket bucket, token token) *Module {
	return &Module{
		repo:      repo,
		bucket:    bucket,
		messanger: messanger,
		token:     token,
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

type bucket interface {
	GeneratePreloadLinkForPlaceClassMedia(
		ctx context.Context,
		placeClassID, sessionID uuid.UUID,
	) (models.PlaceClassUploadMediaLinks, error)

	AcceptUpdatePlaceClassMedia(
		ctx context.Context,
		placeClassID, sessionID uuid.UUID,
	) (models.PlaceMedia, error)

	CancelUpdatePlaceIcon(
		ctx context.Context,
		placeID, sessionID uuid.UUID,
	) error

	CancelUpdateClassIconBanner(
		ctx context.Context,
		placeID, sessionID uuid.UUID,
	) error

	CleanPlaceClassMediaSession(
		ctx context.Context,
		placeID, sessionID uuid.UUID,
	) error
}

type token interface {
	NewUploadPlaceClassMediaToken(
		OwnerAccountID uuid.UUID,
		ResourceID uuid.UUID,
		UploadSessionID uuid.UUID,
	) (string, error)
}
