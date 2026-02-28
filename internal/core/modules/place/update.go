package place

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/models"
)

type UpdateParams struct {
	ClassID uuid.UUID `json:"class_id"`
	Name    string    `json:"name"`
	Address string    `json:"address"`

	Description *string `json:"description"`
	Website     *string `json:"website"`
	Phone       *string `json:"phone"`

	IconKey   *string `json:"icon_key,omitempty"`
	BannerKey *string `json:"banner_key,omitempty"`
}

func (m *Module) Update(
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

	class, err := m.repo.GetPlaceClass(ctx, params.ClassID)
	if err != nil {
		return models.Place{}, err
	}
	if class.DeprecatedAt != nil {
		return models.Place{}, errx.ErrorPlaceClassIsDeprecated.Raise(
			fmt.Errorf("place class %s is deprecated", params.ClassID),
		)
	}

	_, err = m.repo.GetOrgMemberByAccountID(ctx, actor, place.OrganizationID)
	if err != nil {
		return models.Place{}, err
	}

	upd := params.ClassID != place.ClassID ||
		params.Name != place.Name ||
		params.Address != place.Address ||
		!ptrStrEq(params.Description, place.Description) ||
		!ptrStrEq(params.Website, place.Website) ||
		!ptrStrEq(params.Phone, place.Phone)

	if !ptrStrEq(params.IconKey, place.IconKey) {
		iconKey, err := m.updatePlaceIcon(ctx, place, params)
		if err != nil {
			return models.Place{}, fmt.Errorf("failed to validate place icon: %w", err)
		}
		params.IconKey = iconKey
		upd = true
	}

	if !ptrStrEq(params.BannerKey, place.BannerKey) {
		bannerKey, err := m.updatePlaceBanner(ctx, place, params)
		if err != nil {
			return models.Place{}, fmt.Errorf("failed to validate place banner: %w", err)
		}
		params.BannerKey = bannerKey
		upd = true
	}

	if !upd {
		return place, nil
	}

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

func ptrStrEq(a, b *string) bool {
	return (a == nil && b == nil) || (a != nil && b != nil && *a == *b)
}
