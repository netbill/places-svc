package place

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/models"
)

func (m *Module) UpdateStatus(
	ctx context.Context,
	initiator models.Initiator,
	placeID uuid.UUID,
	status string,
) (place models.Place, err error) {
	place, err = m.repo.GetPlaceByID(ctx, placeID)
	if err != nil {
		return models.Place{}, err
	}

	if place.Status == status {
		return place, nil
	}

	if place.OrganizationID != nil {
		if err = m.chekPermissionForUpdatePlace(ctx, initiator, *place.OrganizationID); err != nil {
			return models.Place{}, err
		}
	}

	err = m.repo.Transaction(ctx, func(txCtx context.Context) error {
		place, err = m.repo.UpdatePlaceStatus(txCtx, placeID, status)
		if err != nil {
			return err
		}

		if err = m.messenger.PublishUpdatePlaceStatus(txCtx, place); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return models.Place{}, err
	}

	return place, nil
}
