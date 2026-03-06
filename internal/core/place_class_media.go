package core

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/models"
)

func (m *PlaceClassModule) CreateUploadMediaLinks(
	ctx context.Context,
	placeClassID uuid.UUID,
) (models.PlaceClass, models.UploadPlaceClassMediaLinks, error) {
	class, err := m.repo.Get(ctx, placeClassID)
	if err != nil {
		return models.PlaceClass{}, models.UploadPlaceClassMediaLinks{}, err
	}

	links, err := m.media.CreatePlaceClassIconUploadMediaLinks(ctx, class.ID)
	if err != nil {
		return models.PlaceClass{}, models.UploadPlaceClassMediaLinks{}, err
	}

	return class, models.UploadPlaceClassMediaLinks{
		Icon: links,
	}, nil
}

type DeleteUploadPlaceClassMediaParams struct {
	Icon *string
}

func (m *PlaceClassModule) DeleteUploadMedia(
	ctx context.Context,
	classID uuid.UUID,
	params DeleteUploadPlaceClassMediaParams,
) error {
	_, err := m.repo.Get(ctx, classID)
	if err != nil {
		return err
	}

	if params.Icon != nil {
		if err = m.media.DeleteUploadPlaceClassIcon(ctx, classID, *params.Icon); err != nil {
			return err
		}
	}

	return nil
}
