package place

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/models"
)

type UpdateParams struct {
	ClassID     *uuid.UUID
	Address     *string
	Name        *string
	Description *string
	Icon        *string
	Banner      *string
	Website     *string
	Phone       *string
}

func (s Service) UpdatePlace(
	ctx context.Context,
	initiatorID, placeID uuid.UUID,
	params UpdateParams,
) (place models.Place, err error) {
	place, err = s.repo.GetPlaceByID(ctx, placeID)
	if err != nil {
		return models.Place{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get place by id %s: %w", placeID, err),
		)
	}
	if place.IsNil() {
		return models.Place{}, errx.ErrorPlaceNotFound.Raise(
			fmt.Errorf("place %s not found", placeID),
		)
	}

	if place.OrganizationID != nil {
		if err = s.chekPermissionForManagePlace(ctx, initiatorID, *place.OrganizationID); err != nil {
			return models.Place{}, err
		}
	}

	if params.ClassID != nil {
		classExists, err := s.repo.CheckPlaceClassExists(ctx, *params.ClassID)
		if err != nil {
			return models.Place{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to check place class exists: %w", err),
			)
		}
		if !classExists {
			return models.Place{}, errx.ErrorPlaceClassNotFound.Raise(
				fmt.Errorf("place class %v not found", params.ClassID),
			)
		}
	}
	err = s.repo.Transaction(ctx, func(txCtx context.Context) error {
		place, err = s.repo.UpdatePlaceByID(ctx, placeID, params)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to update place %s: %w", placeID, err),
			)
		}

		err = s.messanger.PublishUpdatePlace(ctx, place)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to publish update place %s: %w", placeID, err),
			)
		}

		return nil
	})

	return place, nil
}

func (s Service) UpdatePlaceStatus(
	ctx context.Context,
	initiatorID, placeID uuid.UUID,
	status string,
) (place models.Place, err error) {
	place, err = s.repo.GetPlaceByID(ctx, placeID)
	if err != nil {
		return models.Place{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get place by id %s: %w", placeID, err),
		)
	}
	if place.IsNil() {
		return models.Place{}, errx.ErrorPlaceNotFound.Raise(
			fmt.Errorf("place %s not found", placeID),
		)
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
		err = s.chekPermissionForManagePlace(ctx, initiatorID, *place.OrganizationID)
		if err != nil {
			return models.Place{}, errx.ErrorNotEnoughRights.Raise(
				fmt.Errorf("initiator %s has no rights to manage place %s", initiatorID, placeID),
			)
		}
	}

	err = s.repo.Transaction(ctx, func(txCtx context.Context) error {
		place, err = s.repo.UpdatePlaceStatus(ctx, placeID, status)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to update place %s status: %w", placeID, err),
			)
		}

		err = s.messanger.PublishUpdatePlaceStatus(ctx, place)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to publish update place %s status: %w", placeID, err),
			)
		}

		return nil
	})

	return place, nil
}

func (s Service) UpdatePlaceVerified(
	ctx context.Context,
	placeID uuid.UUID,
	verified bool,
) (place models.Place, err error) {
	place, err = s.repo.GetPlaceByID(ctx, placeID)
	if err != nil {
		return models.Place{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get place by id %s: %w", placeID, err),
		)
	}
	if place.IsNil() {
		return models.Place{}, errx.ErrorPlaceNotFound.Raise(
			fmt.Errorf("place %s not found", placeID),
		)
	}
	if place.Verified == verified {
		return place, nil
	}

	err = s.repo.Transaction(ctx, func(txCtx context.Context) error {
		place, err = s.repo.UpdatePlaceVerified(ctx, placeID, verified)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to update place %s verified: %w", placeID, err),
			)
		}

		err = s.messanger.PublishUpdatePlaceVerified(ctx, place)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to publish update place %s verified: %w", placeID, err),
			)
		}

		return nil
	})

	return place, nil
}
