package places

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/models"
)

type media interface {
	CreatePlaceIconUploadMediaLinks(
		ctx context.Context,
		placeID uuid.UUID,
	) (models.UploadMediaLink, error)

	UpdatePlaceIcon(
		ctx context.Context,
		orgID uuid.UUID,
		tempKey string,
	) (newKey string, err error)

	DeletePlaceIcon(
		ctx context.Context,
		organizationID uuid.UUID,
		key string,
	) error

	DeleteUploadPlaceIcon(
		ctx context.Context,
		orgID uuid.UUID,
		key string,
	) error

	CreatePlaceBannerUploadMediaLinks(
		ctx context.Context,
		placeID uuid.UUID,
	) (models.UploadMediaLink, error)

	UpdatePlaceBanner(
		ctx context.Context,
		placeID uuid.UUID,
		key string,
	) (string, error)

	DeletePlaceBanner(
		ctx context.Context,
		placeID uuid.UUID,
		key string,
	) error
	DeleteUploadPlaceBanner(
		ctx context.Context,
		placeID uuid.UUID,
		key string,
	) error
}

func (s *Service) CreateUploadMediaLinks(
	ctx context.Context,
	actor models.AccountActor,
	placeID uuid.UUID,
) (models.Place, models.UploadPlaceMediaLinks, error) {
	place, err := s.place.Get(ctx, placeID)
	if err != nil {
		return models.Place{}, models.UploadPlaceMediaLinks{}, err
	}

	_, err = s.org.ValidateOrg(ctx, place.OrganizationID)
	if err != nil {
		return models.Place{}, models.UploadPlaceMediaLinks{}, err
	}

	_, err = s.org.AuthorizeOrgMember(ctx, actor, place.OrganizationID)
	if err != nil {
		return models.Place{}, models.UploadPlaceMediaLinks{}, err
	}

	icon, err := s.media.CreatePlaceIconUploadMediaLinks(ctx, place.ID)
	if err != nil {
		return models.Place{}, models.UploadPlaceMediaLinks{}, err
	}

	banner, err := s.media.CreatePlaceBannerUploadMediaLinks(ctx, place.ID)
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

func (s *Service) DeleteUploadPlaceMedia(
	ctx context.Context,
	actor models.AccountActor,
	placeID uuid.UUID,
	params DeleteUploadPlaceMediaParams,
) error {
	place, err := s.place.Get(ctx, placeID)
	if err != nil {
		return err
	}

	_, err = s.org.AuthorizeOrgHead(ctx, actor, place.OrganizationID)
	if err != nil {
		return err
	}

	_, err = s.org.ValidateOrg(ctx, place.OrganizationID)
	if err != nil {
		return err
	}

	if params.Icon != nil {
		return s.media.DeleteUploadPlaceIcon(ctx, actor, *params.Icon)
	}
	if params.Banner != nil {
		return s.media.DeleteUploadPlaceBanner(ctx, actor, *params.Banner)
	}

	return nil
}
