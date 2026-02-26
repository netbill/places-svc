package controller

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/models"
	"github.com/netbill/places-svc/internal/core/modules/pclass"
	"github.com/netbill/places-svc/internal/core/modules/place"
	"github.com/netbill/restkit/pagi"
)

type placeSvc interface {
	Create(
		ctx context.Context,
		actor models.AccountActor,
		params place.CreateParams,
	) (place models.Place, err error)

	Delete(
		ctx context.Context,
		actor models.AccountActor,
		placeID uuid.UUID,
	) error

	Get(ctx context.Context, placeID uuid.UUID) (models.Place, error)

	GetList(
		ctx context.Context,
		params place.FilterParams,
		limit, offset uint,
	) (pagi.Page[[]models.Place], error)

	Update(
		ctx context.Context,
		actor models.AccountActor,
		placeID uuid.UUID,
		params place.UpdateParams,
	) (place models.Place, err error)

	UpdateStatus(
		ctx context.Context,
		actor models.AccountActor,
		placeID uuid.UUID,
		status string,
	) (place models.Place, err error)

	UpdateVerified(
		ctx context.Context,
		placeID uuid.UUID,
		verified bool,
	) (place models.Place, err error)

	CreateUploadMediaLinks(
		ctx context.Context,
		placeID uuid.UUID,
	) (models.Place, models.UploadPlaceMediaLinks, error)
	DeleteUploadPlaceBanner(
		ctx context.Context,
		actor models.AccountActor,
		placeID uuid.UUID,
		key string,
	) error
	DeleteUploadPlaceIcon(
		ctx context.Context,
		actor models.AccountActor,
		placeID uuid.UUID,
		key string,
	) error
}

type organizationSvc interface {
	Get(ctx context.Context, id uuid.UUID) (models.Organization, error)
}

type placeClassSvc interface {
	Create(
		ctx context.Context,
		params pclass.CreateParams,
	) (class models.PlaceClass, err error)

	Get(ctx context.Context, id uuid.UUID) (models.PlaceClass, error)

	GetList(
		ctx context.Context,
		params pclass.FilterParams,
		limit, offset uint,
	) (pagi.Page[[]models.PlaceClass], error)

	Deprecate(
		ctx context.Context,
		classID uuid.UUID,
	) (models.PlaceClass, error)

	Update(
		ctx context.Context,
		classID uuid.UUID,
		params pclass.UpdateParams,
	) (class models.PlaceClass, err error)

	CreateUploadMediaLinks(
		ctx context.Context,
		placeClassID uuid.UUID,
	) (models.PlaceClass, models.UploadPlaceClassMediaLinks, error)

	DeleteUploadPlaceClassIcon(
		ctx context.Context,
		classID uuid.UUID,
		key string,
	) error
}

type responser interface {
	Render(w http.ResponseWriter, status int, res interface{})
	RenderErr(w http.ResponseWriter, errs ...error)
}

type Modules struct {
	Place placeSvc
	Class placeClassSvc
	Org   organizationSvc
}

type Controller struct {
	modules   *Modules
	responser responser
}

func New(modules *Modules, responser responser) *Controller {
	return &Controller{
		modules:   modules,
		responser: responser,
	}
}
