package place

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/models"
)

func (m *Module) UpdateVerified(
	ctx context.Context,
	placeID uuid.UUID,
	verified bool,
) (place models.Place, err error) {
	place, err = m.repo.GetPlaceByID(ctx, placeID)
	if err != nil {
		return models.Place{}, err
	}

	if place.Verified == verified {
		return place, nil
	}

	err = m.repo.Transaction(ctx, func(txCtx context.Context) error {
		place, err = m.repo.UpdatePlaceVerified(ctx, placeID, verified)
		if err != nil {
			return err
		}

		err = m.messenger.PublishUpdatePlaceVerified(ctx, place)
		if err != nil {
			return err
		}

		return nil
	})

	return place, nil
}
