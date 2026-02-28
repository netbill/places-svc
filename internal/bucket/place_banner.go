package bucket

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/google/uuid"
	"github.com/netbill/awsx"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/models"
)

func CreateTempPlaceBannerKey(placeID uuid.UUID) string {
	return fmt.Sprintf("place/banner/%s/temp/%s", placeID, uuid.New().String())
}

func CreatePlaceBannerKey(placeID uuid.UUID) string {
	return fmt.Sprintf("place/banner/%s/%s", placeID, uuid.New().String())
}

func (s *Storage) CreatePlaceBannerUploadMediaLinks(
	ctx context.Context,
	placeID uuid.UUID,
) (models.UploadMediaLink, error) {
	key := CreateTempPlaceBannerKey(placeID)

	uploadURL, getURL, err := s.s3.PresignPut(
		ctx,
		key,
		s.config.LinkTTL,
	)
	if err != nil {
		return models.UploadMediaLink{}, fmt.Errorf("presigning put for place place banner: %w", err)
	}

	return models.UploadMediaLink{
		Key:        key,
		PreloadUrl: getURL,
		UploadURL:  uploadURL,
	}, nil
}

func (s *Storage) ValidatePlaceBanner(
	ctx context.Context,
	placeID uuid.UUID,
	key string,
) error {
	err := validateTempPlaceBannerKey(placeID, key)
	if err != nil {
		return err
	}

	out, err := s.s3.GetObjectRange(ctx, key, 64*1024)
	switch {
	case errors.Is(err, awsx.ErrNotFound):
		return errx.ErrorNoContentUploaded.Raise(
			fmt.Errorf("place place banner not found for key: %s", key),
		)
	case err != nil:
		return fmt.Errorf("get object range for place place banner: %w", err)
	}
	defer out.Body.Close()

	if err = s.config.PlaceBanner.Validate(out); err != nil {
		switch {
		case errors.Is(err, awsx.ErrorNoContentUploaded):
			return errx.ErrorNoContentUploaded.Raise(err)
		case errors.Is(err, awsx.ErrorSizeExceedsMax):
			return errx.ErrorPlaceBannerContentIsExceedsMax.Raise(err)
		case errors.Is(err, awsx.ErrorResolutionIsInvalid):
			return errx.ErrorPlaceBannerResolutionIsInvalid.Raise(err)
		case errors.Is(err, awsx.ErrorFormatNotAllowed):
			return errx.ErrorPlaceBannerFormatIsNotAllowed.Raise(err)
		default:
			return fmt.Errorf("validate place banner content: %w", err)
		}
	}

	return nil
}

func (s *Storage) DeleteUploadPlaceBanner(
	ctx context.Context,
	placeID uuid.UUID,
	key string,
) error {
	if err := validateTempPlaceBannerKey(placeID, key); err != nil {
		return err
	}

	if err := s.s3.DeleteObject(ctx, key); err != nil {
		return fmt.Errorf("deleting temp place place banner object: %w", err)
	}

	return nil
}

func (s *Storage) DeletePlaceBanner(
	ctx context.Context,
	placeID uuid.UUID,
	key string,
) error {
	if err := validateFinalPlaceBannerKey(placeID, key); err != nil {
		return err
	}

	if err := s.s3.DeleteObject(ctx, key); err != nil {
		return fmt.Errorf("deleting place place banner object: %w", err)
	}

	return nil
}

func (s *Storage) UpdatePlaceBanner(
	ctx context.Context,
	placeID uuid.UUID,
	key string,
) (string, error) {
	if err := validateTempPlaceBannerKey(placeID, key); err != nil {
		return "", err
	}

	finalKey := CreatePlaceBannerKey(placeID)

	if err := s.s3.CopyObject(ctx, key, finalKey); err != nil {
		return "", fmt.Errorf("copying object for place banner: %w", err)
	}

	return finalKey, nil
}

var (
	tempPlaceBannerKeyRe = regexp.MustCompile(
		`^place/banner/([0-9a-fA-F-]{36})/temp/([0-9a-fA-F-]{36})$`,
	)

	finalPlaceBannerKeyRe = regexp.MustCompile(
		`^place/banner/([0-9a-fA-F-]{36})/([0-9a-fA-F-]{36})$`,
	)
)

func validateTempPlaceBannerKey(placeID uuid.UUID, key string) error {
	if key == "" {
		return errx.ErrorPlaceBannerKeyIsInvalid.Raise(fmt.Errorf("empty key"))
	}

	matches := tempPlaceBannerKeyRe.FindStringSubmatch(key)
	if matches == nil {
		return errx.ErrorPlaceBannerKeyIsInvalid.Raise(fmt.Errorf("key %s does not match temp place place banner key pattern", key))
	}

	if matches[1] != placeID.String() {
		return errx.ErrorPlaceBannerKeyIsInvalid.Raise(fmt.Errorf("key %s does not belong to place place %s", key, placeID))
	}

	return nil
}

func validateFinalPlaceBannerKey(placeID uuid.UUID, key string) error {
	if key == "" {
		return errx.ErrorPlaceBannerKeyIsInvalid.Raise(fmt.Errorf("empty key"))
	}

	matches := finalPlaceBannerKeyRe.FindStringSubmatch(key)
	if matches == nil {
		return errx.ErrorPlaceBannerKeyIsInvalid.Raise(fmt.Errorf("key %s does not match final place place banner key pattern", key))
	}

	if matches[1] != placeID.String() {
		return errx.ErrorPlaceBannerKeyIsInvalid.Raise(fmt.Errorf("key %s does not belong to place place %s", key, placeID))
	}

	return nil
}
