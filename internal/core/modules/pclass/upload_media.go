package pclass

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/models"
)

func (m *Module) CreateUploadMediaLinks(
	ctx context.Context,
	placeClassID uuid.UUID,
) (models.PlaceClass, models.UploadPlaceClassMediaLinks, error) {
	class, err := m.repo.GetPlaceClass(ctx, placeClassID)
	if err != nil {
		return models.PlaceClass{}, models.UploadPlaceClassMediaLinks{}, err
	}

	links, err := m.bucket.CreatePlaceClassIconUploadMediaLinks(ctx, class.ID)
	if err != nil {
		return models.PlaceClass{}, models.UploadPlaceClassMediaLinks{}, err
	}

	return class, models.UploadPlaceClassMediaLinks{
		Icon: links,
	}, nil
}

func (m *Module) DeleteUploadPlaceClassIcon(
	ctx context.Context,
	classID uuid.UUID,
	key string,
) error {
	_, err := m.repo.GetPlaceClass(ctx, classID)
	if err != nil {
		return err
	}

	return m.bucket.DeleteUploadPlaceClassIcon(ctx, classID, key)
}

func (m *Module) updatePlaceClassIcon(
	ctx context.Context,
	class models.PlaceClass,
	params UpdateParams,
) (newKey *string, err error) {
	if params.IconKey != nil {
		err = m.bucket.ValidatePlaceClassIcon(ctx, class.ID, *params.IconKey)
		if err != nil {
			return nil, fmt.Errorf("failed to validate place class icon: %w", err)
		}

		iconKey, err := m.bucket.UpdatePlaceClassIcon(ctx, class.ID, *params.IconKey)
		if err != nil {
			return nil, fmt.Errorf("failed to update place class icon: %w", err)
		}

		newKey = &iconKey
	}

	if class.IconKey != nil && params.IconKey != nil || (params.IconKey == nil && class.IconKey != nil) {
		err = m.bucket.DeletePlaceClassIcon(ctx, class.ID, *class.IconKey)
		if err != nil {
			return nil, fmt.Errorf("failed to delete place class icon: %w", err)
		}
	}

	return newKey, nil
}
