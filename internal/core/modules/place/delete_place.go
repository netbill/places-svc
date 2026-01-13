package place

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
)

func (s Service) DeletePlace(
	ctx context.Context,
	initiatorID, placeID uuid.UUID,
) error {
	place, err := s.repo.GetPlaceByID(ctx, placeID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get place by id %s: %w", placeID, err),
		)
	}
	if place.IsNil() {
		return errx.ErrorPlaceNotFound.Raise(
			fmt.Errorf("place %s not found", placeID),
		)
	}

	if place.OrganizationID != nil {
		err = s.chekPermissionForManagePlace(ctx, initiatorID, *place.OrganizationID)
		if err != nil {
			return errx.ErrorNotEnoughRights.Raise(
				fmt.Errorf("initiator %s has no rights to manage place %s", initiatorID, placeID),
			)
		}
	}

	return s.repo.Transaction(ctx, func(txCtx context.Context) error {
		err = s.repo.DeletePlaceByID(ctx, placeID)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to delete place %s: %w", placeID, err),
			)
		}

		err = s.messanger.PublishDeletePlace(ctx, placeID)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to publish delete place %s: %w", placeID, err),
			)
		}

		return nil
	})
}
