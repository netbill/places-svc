package places

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/models"
	"github.com/netbill/restkit/pagi"
)

type placeRepo interface {
	Create(ctx context.Context, params CreateParams) (models.Place, error)

	Get(ctx context.Context, placeID uuid.UUID) (models.Place, error)
	GetListByIDs(ctx context.Context, placeIDs []uuid.UUID) ([]models.Place, error)
	GetList(ctx context.Context, params FilterPlaceParams, limit, offset uint) (pagi.Page[[]models.Place], error)

	Update(ctx context.Context, placeID uuid.UUID, params UpdateParams) (models.Place, error)
	UpdateStatus(ctx context.Context, placeID uuid.UUID, status string) (models.Place, error)
	UpdateVerified(ctx context.Context, placeID uuid.UUID, verified bool) (models.Place, error)
	UpdateClass(ctx context.Context, placeID uuid.UUID, classID uuid.UUID) (models.Place, error)

	Delete(ctx context.Context, id uuid.UUID) error
}

type classRepo interface {
	Get(ctx context.Context, classID uuid.UUID) (models.PlaceClass, error)
	Exists(ctx context.Context, classID uuid.UUID) (bool, error)
}

type tombstoneRepo interface {
	BuryPlace(ctx context.Context, placeID uuid.UUID) error
	PlaceIsBuried(ctx context.Context, placeID uuid.UUID) (bool, error)
}

type transaction interface {
	Transaction(ctx context.Context, fn func(context.Context) error) error
}
