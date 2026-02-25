package place

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/models"
)

type UpdateParams struct {
	Address string `json:"address"`
	Name    string `json:"name"`

	Description *string `json:"description"`
	Website     *string `json:"website"`
	Phone       *string `json:"phone"`

	IconKey   *string `json:"icon_key,omitempty"`
	BannerKey *string `json:"banner_key,omitempty"`
}

func (m *Module) ConfirmUpdateSession(
	ctx context.Context,
	actor models.AccountActor,
	placeID uuid.UUID,
	params UpdateParams,
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

	_, err = m.repo.GetOrgMemberByAccountID(ctx, actor, place.OrganizationID)
	if err != nil {
		return models.Place{}, err
	}

	icon, err := m.updatePlaceIcon(ctx, place, params)
	if err != nil {
		return models.Place{}, err
	}
	params.IconKey = icon

	banner, err := m.updatePlaceBanner(ctx, place, params)
	if err != nil {
		return models.Place{}, err
	}
	params.BannerKey = banner

	err = m.repo.Transaction(ctx, func(txCtx context.Context) error {
		place, err = m.repo.UpdatePlaceByID(txCtx, placeID, params)
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
