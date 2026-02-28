package place

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/models"
)

func (m *Module) Delete(
	ctx context.Context,
	actor models.AccountActor,
	placeID uuid.UUID,
) error {
	place, err := m.repo.GetPlaceByID(ctx, placeID)
	if err != nil {
		return err
	}

	org, err := m.repo.GetOrganization(ctx, place.OrganizationID)
	if err != nil {
		return err
	}
	if org.Status == models.OrganizationStatusSuspended {
		return errx.ErrorOrganizationIsSuspended.Raise(
			fmt.Errorf("organization %s is suspended", place.OrganizationID),
		)
	}

	member, err := m.repo.GetOrgMemberByAccountID(ctx, actor, place.OrganizationID)
	if err != nil {
		return err
	}
	if !member.Head {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("only organization head can delete places"),
		)
	}

	buried, err := m.repo.PlaceIsBuried(ctx, placeID)
	if err != nil {
		return err
	}
	if buried {
		return errx.ErrorPlaceDeleted.Raise(
			fmt.Errorf("place with id %s is already deleted", placeID),
		)
	}

	return m.repo.Transaction(ctx, func(ctx context.Context) error {
		if err = m.repo.BuryPlace(ctx, placeID); err != nil {
			return err
		}

		if err = m.repo.DeletePlaceByID(ctx, placeID); err != nil {
			return err
		}

		return m.messenger.PublishDeletePlace(ctx, placeID)
	})
}
