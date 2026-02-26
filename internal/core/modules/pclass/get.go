package pclass

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/models"
	"github.com/netbill/restkit/pagi"
)

func (m *Module) Get(ctx context.Context, id uuid.UUID) (models.PlaceClass, error) {
	class, err := m.repo.GetPlaceClass(ctx, id)
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

func (m *Module) GetList(
	ctx context.Context,
	params FilterParams,
	limit, offset uint,
) (pagi.Page[[]models.PlaceClass], error) {
	classes, err := m.repo.GetPlaceClasses(ctx, params, limit, offset)
	if err != nil {
		return pagi.Page[[]models.PlaceClass]{}, err
	}

	return classes, nil
}

func (m *Module) GetByIDs(
	ctx context.Context,
	ids []uuid.UUID,
) ([]models.PlaceClass, error) {
	res, err := m.repo.GetPlaceClassesByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	return res, nil
}
