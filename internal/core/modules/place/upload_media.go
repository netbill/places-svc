package place

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/models"
)

func (m *Module) CreateUploadMediaLinks(
	ctx context.Context,
	placeID uuid.UUID,
) (models.Place, models.UploadPlaceMediaLinks, error) {
	place, err := m.repo.GetPlaceByID(ctx, placeID)
	if err != nil {
		return models.Place{}, models.UploadPlaceMediaLinks{}, err
	}

	links, err := m.bucket.CreatePlaceIconUploadMediaLinks(ctx, place.ID)
	if err != nil {
		return models.Place{}, models.UploadPlaceMediaLinks{}, err
	}

	return place, models.UploadPlaceMediaLinks{
		Icon: links,
	}, nil
}

func (m *Module) DeleteUploadPlaceIcon(
	ctx context.Context,
	actor models.AccountActor,
	placeID uuid.UUID,
	key string,
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

	_, err = m.repo.GetOrgMemberByAccountID(ctx, actor, place.OrganizationID)
	if err != nil {
		return err
	}

	err = m.bucket.DeleteUploadPlaceIcon(ctx, actor, key)
	if err != nil {
		return err
	}

	return nil
}

func (m *Module) updatePlaceIcon(
	ctx context.Context,
	place models.Place,
	params UpdateParams,
) (newKey *string, err error) {
	if params.IconKey != nil {
		if err = m.bucket.ValidatePlaceIcon(ctx, place.ID, *params.IconKey); err != nil {
			return nil, fmt.Errorf("failed to validate place icon: %w", err)
		}

		iconKey, err := m.bucket.UpdatePlaceIcon(ctx, place.ID, *params.IconKey)
		if err != nil {
			return nil, fmt.Errorf("failed to update place icon: %w", err)
		}

		if err = m.bucket.DeletePlaceIcon(ctx, place.ID, iconKey); err != nil {
			return nil, fmt.Errorf("failed to delete place icon: %w", err)
		}

		newKey = &iconKey
	}

	if place.IconKey != nil {
		if err = m.bucket.DeletePlaceIcon(ctx, place.ID, *place.IconKey); err != nil {
			return nil, fmt.Errorf("failed to delete place icon: %w", err)
		}
	}

	return newKey, nil
}

func (m *Module) DeleteUploadPlaceBanner(
	ctx context.Context,
	actor models.AccountActor,
	placeID uuid.UUID,
	key string,
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

	_, err = m.repo.GetOrgMemberByAccountID(ctx, actor, place.OrganizationID)
	if err != nil {
		return err
	}

	err = m.bucket.DeleteUploadPlaceBanner(ctx, actor, key)
	if err != nil {
		return err
	}

	return nil
}

func (m *Module) updatePlaceBanner(
	ctx context.Context,
	place models.Place,
	params UpdateParams,
) (newKey *string, err error) {
	if params.BannerKey != nil {
		if err = m.bucket.ValidatePlaceBanner(ctx, place.ID, *params.BannerKey); err != nil {
			return nil, fmt.Errorf("failed to validate place banner: %w", err)
		}

		key, err := m.bucket.UpdatePlaceBanner(ctx, place.ID, *params.BannerKey)
		if err != nil {
			return nil, fmt.Errorf("failed to update place banner: %w", err)
		}

		if err = m.bucket.DeletePlaceBanner(ctx, place.ID, key); err != nil {
			return nil, fmt.Errorf("failed to delete place banner: %w", err)
		}

		newKey = &key
	}

	if place.BannerKey != nil {
		if err = m.bucket.DeletePlaceBanner(ctx, place.ID, *place.BannerKey); err != nil {
			return nil, fmt.Errorf("failed to delete place banner: %w", err)
		}
	}

	return newKey, nil
}
