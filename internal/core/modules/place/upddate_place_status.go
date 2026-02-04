package place

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/models"
)

func (m *Module) UpdatePlaceStatus(
	ctx context.Context,
	initiator models.InitiatorData,
	placeID uuid.UUID,
	status string,
) (place models.Place, err error) {
	place, err = m.repo.GetPlaceByID(ctx, placeID)
	if err != nil {
		return models.Place{}, err
	}

	if place.Status == models.OrganizationStatusSuspended {
		return models.Place{}, errx.ErrorPlaceStatusSuspended.Raise(
			fmt.Errorf("place %s is suspended and status cannot be changed", placeID),
		)
	}
	if status == models.OrganizationStatusSuspended {
		return models.Place{}, errx.ErrorCannotSetPlaceStatusSuspend.Raise(
			fmt.Errorf("place %s status cannot be changed to suspended", placeID),
		)
	}
	if place.Status == status {
		return place, nil
	}

	if place.OrganizationID != nil {
		err = m.chekPermissionForManagePlace(ctx, initiator, *place.OrganizationID)
		if err != nil {
			return models.Place{}, err
		}
	}

	err = m.repo.Transaction(ctx, func(txCtx context.Context) error {
		place, err = m.repo.UpdatePlaceStatus(ctx, placeID, status)
		if err != nil {
			return err
		}

		err = m.messanger.PublishUpdatePlaceStatus(ctx, place)
		if err != nil {
			return err
		}

		return nil
	})

	return place, nil
}
