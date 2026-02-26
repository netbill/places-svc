package place

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/models"
	"github.com/netbill/restkit/pagi"
)

type Module struct {
	repo      repo
	bucket    bucket
	messenger messenger
	territory checkerTerritory
}

func New(
	repo repo,
	bucket bucket,
	messenger messenger,
	territory checkerTerritory,
) *Module {
	return &Module{
		repo:      repo,
		bucket:    bucket,
		messenger: messenger,
		territory: territory,
	}
}

type repo interface {
	CreatePlace(ctx context.Context, params CreateParams) (models.Place, error)

	GetPlaceByID(ctx context.Context, id uuid.UUID) (models.Place, error)
	GetPlacesByIDs(ctx context.Context, ids []uuid.UUID) ([]models.Place, error)
	GetPlaces(ctx context.Context, params FilterParams, limit, offset uint) (pagi.Page[[]models.Place], error)

	UpdatePlaceByID(ctx context.Context, id uuid.UUID, params UpdateParams) (models.Place, error)
	UpdatePlaceStatus(ctx context.Context, placeID uuid.UUID, status string) (models.Place, error)
	UpdatePlaceVerified(ctx context.Context, placeID uuid.UUID, verified bool) (models.Place, error)
	UpdateClassForPlace(ctx context.Context, placeID uuid.UUID, classID uuid.UUID) (models.Place, error)

	DeletePlaceByID(ctx context.Context, id uuid.UUID) error

	GetPlaceClass(ctx context.Context, classID uuid.UUID) (models.PlaceClass, error)

	GetOrganization(ctx context.Context, id uuid.UUID) (models.Organization, error)
	GetOrgMemberByAccountID(ctx context.Context, organizationID, accountID uuid.UUID) (models.OrgMember, error)

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type messenger interface {
	PublishCreatePlace(ctx context.Context, place models.Place) error
	PublishUpdatePlace(ctx context.Context, place models.Place) error
	PublishDeletePlace(ctx context.Context, placeID uuid.UUID) error
}

type bucket interface {
	CreatePlaceIconUploadMediaLinks(
		ctx context.Context,
		classID uuid.UUID,
	) (models.UploadMediaLink, error)

	ValidatePlaceIcon(
		ctx context.Context,
		classID uuid.UUID,
		key string,
	) error

	DeleteUploadPlaceIcon(
		ctx context.Context,
		classID uuid.UUID,
		key string,
	) error

	DeletePlaceIcon(
		ctx context.Context,
		classID uuid.UUID,
		key string,
	) error

	UpdatePlaceIcon(
		ctx context.Context,
		classID uuid.UUID,
		key string,
	) (string, error)

	CreatePlaceBannerUploadMediaLinks(
		ctx context.Context,
		classID uuid.UUID,
	) (models.UploadMediaLink, error)

	ValidatePlaceBanner(
		ctx context.Context,
		classID uuid.UUID,
		key string,
	) error

	DeleteUploadPlaceBanner(
		ctx context.Context,
		classID uuid.UUID,
		key string,
	) error

	DeletePlaceBanner(
		ctx context.Context,
		classID uuid.UUID,
		key string,
	) error

	UpdatePlaceBanner(
		ctx context.Context,
		classID uuid.UUID,
		key string,
	) (string, error)
}

type checkerTerritory interface {
	ContainsLatLng(lat, lng float64) bool
}
