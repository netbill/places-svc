package place

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/models"
	"github.com/netbill/restkit/pagi"
	"github.com/paulmach/orb"
)

func (m *Module) Get(ctx context.Context, placeID uuid.UUID) (models.Place, error) {
	place, err := m.repo.GetPlaceByID(ctx, placeID)
	if err != nil {
		return models.Place{}, err
	}

	return place, nil
}

type FilterParams struct {
	Class          *FilterClassParams
	Near           *FilterNearParams
	OrganizationID *uuid.UUID
	Status         []string
	Verified       *bool

	BestMatch *string

	Address     *string
	Name        *string
	Description *string

	Website *string
	Phone   *string
}

type FilterClassParams struct {
	ClassID  []uuid.UUID
	Parents  bool
	Children bool
}

type FilterNearParams struct {
	Point   orb.Point
	RadiusM uint
}

func (m *Module) GetList(
	ctx context.Context,
	params FilterParams,
	limit, offset uint,
) (pagi.Page[[]models.Place], error) {
	res, err := m.repo.GetPlaces(ctx, params, limit, offset)
	if err != nil {
		return pagi.Page[[]models.Place]{}, err
	}

	return res, nil
}
