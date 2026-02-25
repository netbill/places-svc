package place

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/models"
)

func (m *Module) UpdateClass(
	ctx context.Context,
	actor models.AccountActor,
	placeID uuid.UUID,
	classID uuid.UUID,
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

	member, err := m.repo.GetOrgMemberByAccountID(ctx, actor, place.OrganizationID)
	if err != nil {
		return models.Place{}, err
	}
	if !member.Head {
		return models.Place{}, errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("only organization head can update places"),
		)
	}

	exist, err := m.repo.CheckPlaceClassExists(ctx, classID)
	if err != nil {
		return models.Place{}, err
	}
	if !exist {
		return models.Place{}, errx.ErrorPlaceClassCodeExists.Raise(
			fmt.Errorf("place class %s does not exist", classID),
		)
	}

	err = m.repo.Transaction(ctx, func(txCtx context.Context) error {
		place, err = m.repo.UpdateClassForPlace(ctx, placeID, classID)
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
