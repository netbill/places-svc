package place

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/models"
)

func (m *Module) DeletePlace(
	ctx context.Context,
	initiator models.InitiatorData,
	placeID uuid.UUID,
) error {
	place, err := m.repo.GetPlaceByID(ctx, placeID)
	if err != nil {
		return err
	}

	if place.OrganizationID != nil {
		err = m.chekPermissionForManagePlace(ctx, initiator, *place.OrganizationID)
		if err != nil {
			return err
		}
	}

	return m.repo.Transaction(ctx, func(txCtx context.Context) error {
		err = m.repo.DeletePlaceByID(ctx, placeID)
		if err != nil {
			return err
		}

		err = m.messanger.PublishDeletePlace(ctx, placeID)
		if err != nil {
			return err
		}

		return nil
	})
}
