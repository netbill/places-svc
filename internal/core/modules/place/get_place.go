package place

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/models"
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
