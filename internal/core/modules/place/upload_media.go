package place

import (
	"context"
	"fmt"

	"github.com/google/uuid"
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
	key string,
) error {
	err := m.bucket.DeleteUploadPlaceIcon(ctx, actor, key)
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
		err = m.bucket.ValidatePlaceIcon(ctx, place.ID, *params.IconKey)
		if err != nil {
			return nil, fmt.Errorf("failed to validate place icon: %w", err)
		}

		iconKey, err := m.bucket.UpdatePlaceIcon(ctx, place.ID, *params.IconKey)
		if err != nil {
			return nil, fmt.Errorf("failed to update place icon: %w", err)
		}

		newKey = &iconKey
	}

	if place.IconKey != nil && params.IconKey != nil || (params.IconKey == nil && place.IconKey != nil) {
		err = m.bucket.DeletePlaceIcon(ctx, place.ID, *place.IconKey)
		if err != nil {
			return nil, fmt.Errorf("failed to delete place icon: %w", err)
		}
	}

	return newKey, nil
}

func (m *Module) DeleteUploadPlaceBanner(
	ctx context.Context,
	actor models.AccountActor,
	key string,
) error {
	err := m.bucket.DeleteUploadPlaceBanner(ctx, actor, key)
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
		err = m.bucket.ValidatePlaceBanner(ctx, place.ID, *params.BannerKey)
		if err != nil {
			return nil, fmt.Errorf("failed to validate place banner: %w", err)
		}

		BannerKey, err := m.bucket.UpdatePlaceBanner(ctx, place.ID, *params.BannerKey)
		if err != nil {
			return nil, fmt.Errorf("failed to update place banner: %w", err)
		}

		newKey = &BannerKey
	}

	if place.BannerKey != nil && params.BannerKey != nil || (params.BannerKey == nil && place.BannerKey != nil) {
		err = m.bucket.DeletePlaceBanner(ctx, place.ID, *place.BannerKey)
		if err != nil {
			return nil, fmt.Errorf("failed to delete place banner: %w", err)
		}
	}

	return newKey, nil
}
