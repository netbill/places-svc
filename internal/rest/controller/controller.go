package controller

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/logium"
	"github.com/netbill/places-svc/internal/core/models"
	"github.com/netbill/places-svc/internal/core/modules/pclass"
	"github.com/netbill/places-svc/internal/core/modules/place"
	"github.com/netbill/restkit/pagi"
)

type placeSvc interface {
	CreatePlace(
		ctx context.Context,
		initiatorID uuid.UUID,
		params place.CreateParams,
	) (place models.Place, err error)

	GetPlace(ctx context.Context, placeID uuid.UUID) (models.Place, error)
	GetPlaces(
		ctx context.Context,
		params place.FilterParams,
		limit, offset uint,
	) (places pagi.Page[[]models.Place], err error)

	UpdatePlace(
		ctx context.Context,
		initiatorID, placeID uuid.UUID,
		params place.UpdateParams,
	) (place models.Place, err error)
	UpdatePlaceStatus(
		ctx context.Context,
		initiatorID, placeID uuid.UUID,
		status string,
	) (place models.Place, err error)
	UpdatePlaceVerified(
		ctx context.Context,
		placeID uuid.UUID,
		verified bool,
	) (place models.Place, err error)

	DeletePlace(ctx context.Context, initiatorID, placeID uuid.UUID) error
}

type placeClassSvc interface {
	CreatePlaceClass(ctx context.Context, params pclass.CreateParams) (class models.PlaceClass, err error)
	GetPlaceClasses(
		ctx context.Context,
		params pclass.FilterParams,
		limit, offset uint,
	) (classes pagi.Page[[]models.PlaceClass], err error)
	GetPlaceClass(ctx context.Context, id uuid.UUID) (models.PlaceClass, error)
	UpdatePlaceClass(
		ctx context.Context,
		ID uuid.UUID,
		params pclass.UpdateParams,
	) (class models.PlaceClass, err error)
	DeletePlaceClass(ctx context.Context, classID uuid.UUID) error

	ReplacePlacesClass(ctx context.Context, oldClassID, newClassID uuid.UUID) error
}

type core struct {
	placeSvc
	placeClassSvc
}

type Controller struct {
	core core
	log  *logium.Logger
}

func New(
	placeSvc placeSvc,
	placeClassSvc placeClassSvc,
	log *logium.Logger,
) Controller {
	return Controller{
		core: core{
			placeSvc:      placeSvc,
			placeClassSvc: placeClassSvc,
		},
		log: log,
	}
}
