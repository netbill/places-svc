package core

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/models"
)

func (m *PlaceModule) CreateUploadMediaLinks(
	ctx context.Context,
	placeID uuid.UUID,
) (models.Place, models.UploadPlaceMediaLinks, error) {
	place, err := m.repo.Get(ctx, placeID)
	if err != nil {
		return models.Place{}, models.UploadPlaceMediaLinks{}, err
	}

	icon, err := m.media.CreatePlaceIconUploadMediaLinks(ctx, place.ID)
	if err != nil {
		return models.Place{}, models.UploadPlaceMediaLinks{}, err
	}

	banner, err := m.media.CreatePlaceBannerUploadMediaLinks(ctx, place.ID)
	if err != nil {
		return models.Place{}, models.UploadPlaceMediaLinks{}, err
	}

	return place, models.UploadPlaceMediaLinks{
		Icon:   icon,
		Banner: banner,
	}, nil
}

type DeleteUploadPlaceMediaParams struct {
	Icon   *string
	Banner *string
}

func (m *PlaceModule) DeleteUploadPlaceMedia(
	ctx context.Context,
	actor models.AccountActor,
	placeID uuid.UUID,
	params DeleteUploadPlaceMediaParams,
) error {
	place, err := m.repo.Get(ctx, placeID)
	if err != nil {
		return err
	}

	_, err = m.auth.authorizeOrgHead(ctx, actor, place.OrganizationID)
	if err != nil {
		return err
	}

	_, err = m.auth.validateOrg(ctx, place.OrganizationID)
	if err != nil {
		return err
	}

	if params.Icon != nil {
		return m.media.DeleteUploadPlaceIcon(ctx, actor, *params.Icon)
	}
	if params.Banner != nil {
		return m.media.DeleteUploadPlaceBanner(ctx, actor, *params.Banner)
	}

	return nil
}
