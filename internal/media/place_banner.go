package media

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

func (s *Uploader) CreatePlaceBannerUploadMediaLinks(
	ctx context.Context,
	placeID uuid.UUID,
) (models.UploadMediaLink, error) {
	key := createTempPlaceBannerKey(placeID)

	uploadURL, getURL, err := s.s3.PresignPut(ctx, key, s.config.LinkTTL)
	if err != nil {
		return models.UploadMediaLink{}, fmt.Errorf("presigning put for place banner: %w", err)
	}

	return models.UploadMediaLink{
		Key:        key,
		PreloadUrl: getURL,
		UploadURL:  uploadURL,
	}, nil
}

func (s *Uploader) UpdatePlaceBanner(
	ctx context.Context,
	placeID uuid.UUID,
	key string,
) (string, error) {
	err := validateTempPlaceBannerKey(placeID, key)
	if err != nil {
		return "", err
	}

	out, err := s.s3.GetObjectRange(ctx, key, 64*1024)
	switch {
	case errors.Is(err, awsx.ErrNotFound):
		return "", errx.ErrorPlaceBannerIsInvalid.Raise(
			fmt.Errorf("place banner not found for key: %s", key),
		)
	case err != nil:
		return "", fmt.Errorf("get object range for place banner: %w", err)
	}
	defer out.Body.Close()

	if err = s.config.PlaceBanner.Validate(out); err != nil {
		return "", errx.ErrorPlaceBannerIsInvalid.Raise(
			fmt.Errorf("validating place banner content for key %s: %w", key, err),
		)
	}

	finalKey := createPlaceBannerKey(placeID)

	if err = s.s3.CopyObject(ctx, key, finalKey); err != nil {
		return "", fmt.Errorf("copying object for place banner: %w", err)
	}

	return finalKey, nil
}

func (s *Uploader) DeleteUploadPlaceBanner(
	ctx context.Context,
	placeID uuid.UUID,
	key string,
) error {
	if err := validateTempPlaceBannerKey(placeID, key); err != nil {
		return err
	}

	if err := s.s3.DeleteObject(ctx, key); err != nil {
		return fmt.Errorf("deleting temp place banner object: %w", err)
	}

	return nil
}

func (s *Uploader) DeletePlaceBanner(
	ctx context.Context,
	placeID uuid.UUID,
	key string,
) error {
	if err := validateFinalPlaceBannerKey(placeID, key); err != nil {
		return err
	}

	if err := s.s3.DeleteObject(ctx, key); err != nil {
		return fmt.Errorf("deleting place banner object: %w", err)
	}

	return nil
}

func createTempPlaceBannerKey(placeID uuid.UUID) string {
	return fmt.Sprintf("place/banner/%s/temp/%s", placeID, uuid.New().String())
}

var tempPlaceBannerKeyRe = regexp.MustCompile(
	`^place/banner/([0-9a-fA-F-]{36})/temp/([0-9a-fA-F-]{36})$`,
)

func validateTempPlaceBannerKey(placeID uuid.UUID, key string) error {
	if key == "" {
		return errx.ErrorPlaceBannerIsInvalid.Raise(fmt.Errorf("empty key"))
	}

	matches := tempPlaceBannerKeyRe.FindStringSubmatch(key)
	if matches == nil {
		return errx.ErrorPlaceBannerIsInvalid.Raise(fmt.Errorf("key %s does not match temp place banner key pattern", key))
	}

	if matches[1] != placeID.String() {
		return errx.ErrorPlaceBannerIsInvalid.Raise(fmt.Errorf("key %s does not belong to place %s", key, placeID))
	}

	return nil
}

func createPlaceBannerKey(placeID uuid.UUID) string {
	return fmt.Sprintf("place/banner/%s/%s", placeID, uuid.New().String())
}

var finalPlaceBannerKeyRe = regexp.MustCompile(
	`^place/banner/([0-9a-fA-F-]{36})/([0-9a-fA-F-]{36})$`,
)

func validateFinalPlaceBannerKey(placeID uuid.UUID, key string) error {
	if key == "" {
		return errx.ErrorPlaceBannerIsInvalid.Raise(fmt.Errorf("empty key"))
	}

	matches := finalPlaceBannerKeyRe.FindStringSubmatch(key)
	if matches == nil {
		return errx.ErrorPlaceBannerIsInvalid.Raise(fmt.Errorf("key %s does not match final place banner key pattern", key))
	}

	if matches[1] != placeID.String() {
		return errx.ErrorPlaceBannerIsInvalid.Raise(fmt.Errorf("key %s does not belong to place %s", key, placeID))
	}

	return nil
}
