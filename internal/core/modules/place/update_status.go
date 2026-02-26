package place

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/models"
)

func (m *Module) UpdateStatus(
	ctx context.Context,
	actor models.AccountActor,
	placeID uuid.UUID,
	status string,
) (place models.Place, err error) {
	place, err = m.repo.GetPlaceByID(ctx, placeID)
	if err != nil {
		return models.Place{}, err
	}

	org, err := m.repo.GetOrganization(ctx, place.OrganizationID)
	if err != nil {
		return models.Place{}, err
	}

	if org.Status == models.OrganizationStatusSuspended {
		return models.Place{}, errx.ErrorOrganizationIsSuspended.Raise(
			fmt.Errorf("organization %s is suspended", place.OrganizationID),
		)
	}

	switch {
	case status == models.PlaceStatusSuspended:
		return models.Place{}, errx.ErrorCannotSetStatusSuspend.Raise(
			fmt.Errorf("cannot set status to %s, suspended", status),
		)
	case status == place.Status:
		return place, nil
	case status != models.PlaceStatusActive, status != models.PlaceStatusInactive:
		return models.Place{}, errx.ErrorPlaceStatusIsInvalid.Raise(
			fmt.Errorf("cannot set status to %s, suspended", status),
		)
	}

	member, err := m.repo.GetOrgMemberByAccountID(ctx, actor, place.OrganizationID)
	if err != nil {
		return models.Place{}, err
	}
	if !member.Head {
		return models.Place{}, errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("only organization head can update status for places"),
		)
	}

	err = m.repo.Transaction(ctx, func(txCtx context.Context) error {
		place, err = m.repo.UpdatePlaceStatus(txCtx, placeID, status)
		if err != nil {
			return err
		}

		if err = m.messenger.PublishUpdatePlace(txCtx, place); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return models.Place{}, err
	}

	return place, nil
}
