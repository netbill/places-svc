package place

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/models"
	"github.com/netbill/restkit/pagi"
	"github.com/paulmach/orb"
)

func (s Service) GetPlace(ctx context.Context, placeID uuid.UUID) (models.Place, error) {
	place, err := s.repo.GetPlaceByID(ctx, placeID)
	if err != nil {
		return models.Place{}, errx.ErrorInternal.Raise(
			fmt.Errorf("place with id %v not found", placeID),
		)
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

func (s Service) GetPlaces(ctx context.Context, params FilterParams, limit, offset uint) (pagi.Page[[]models.Place], error) {
	if limit == 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	res, err := s.repo.GetPlaces(ctx, params, limit, offset)
	if err != nil {
		return pagi.Page[[]models.Place]{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to filter places: %w", err),
		)
	}

	return res, nil
}
