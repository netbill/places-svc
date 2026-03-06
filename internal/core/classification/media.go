package classification

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/models"
)

func (s *Service) CreateUploadMediaLinks(
	ctx context.Context,
	placeClassID uuid.UUID,
) (models.PlaceClass, models.UploadPlaceClassMediaLinks, error) {
	class, err := s.repo.Get(ctx, placeClassID)
	if err != nil {
		return models.PlaceClass{}, models.UploadPlaceClassMediaLinks{}, err
	}

	links, err := s.media.CreatePlaceClassIconUploadMediaLinks(ctx, class.ID)
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

func (s *Service) DeleteUploadMedia(
	ctx context.Context,
	classID uuid.UUID,
	params DeleteUploadPlaceClassMediaParams,
) error {
	_, err := s.repo.Get(ctx, classID)
	if err != nil {
		return err
	}

	if params.Icon != nil {
		if err = s.media.DeleteUploadPlaceClassIcon(ctx, classID, *params.Icon); err != nil {
			return err
		}
	}

	return nil
}
