package classification

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/models"
	"github.com/netbill/restkit/pagi"
)

type repo interface {
	Create(ctx context.Context, params CreateParams) (models.PlaceClass, error)

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
		params UpdateParams,
	) (models.PlaceClass, error)
	Deprecated(
		ctx context.Context,
		classID uuid.UUID,
		value bool,
	) (models.PlaceClass, error)

	CheckParentCycle(ctx context.Context, classID, parentID uuid.UUID) (bool, error)
	CheckHasChildren(ctx context.Context, classID uuid.UUID) (bool, error)
}

type transaction interface {
	Transaction(ctx context.Context, fn func(context.Context) error) error
}
