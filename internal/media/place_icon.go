package media

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/google/uuid"
	"github.com/netbill/awsx"
	"github.com/netbill/places-svc/internal/errx"
	"github.com/netbill/places-svc/internal/models"
)

func (s *Uploader) CreatePlaceIconUploadMediaLinks(
	ctx context.Context,
	placeID uuid.UUID,
) (models.UploadMediaLink, error) {
	key := createTempPlaceIconKey(placeID)

	uploadURL, getURL, err := s.s3.PresignPut(ctx, key, s.config.LinkTTL)
	if err != nil {
		return models.UploadMediaLink{}, fmt.Errorf("presigning put for place icon: %w", err)
	}

	return models.UploadMediaLink{
		Key:        key,
		PreloadUrl: getURL,
		UploadURL:  uploadURL,
	}, nil
}

func (s *Uploader) UpdatePlaceIcon(
	ctx context.Context,
	placeID uuid.UUID,
	key string,
) (string, error) {
	err := validateTempPlaceIconKey(placeID, key)
	if err != nil {
		return "", err
	}

	out, err := s.s3.GetObjectRange(ctx, key, 64*1024)
	switch {
	case errors.Is(err, awsx.ErrNotFound):
		return "", errx.ErrorPlaceIconIsInvalid.Raise(
			fmt.Errorf("place icon not found for key: %s", key),
		)
	case err != nil:
		return "", fmt.Errorf("get object range for place icon: %w", err)
	}
	defer out.Body.Close()

	if err = s.config.PlaceIcon.Validate(out); err != nil {
		return "", errx.ErrorPlaceIconIsInvalid.Raise(
			fmt.Errorf("validating place icon content for key %s: %w", key, err),
		)
	}

	finalKey := createPlaceIconKey(placeID)

	if err = s.s3.CopyObject(ctx, key, finalKey); err != nil {
		return "", fmt.Errorf("copying object for place icon: %w", err)
	}

	return finalKey, nil
}

func (s *Uploader) DeleteUploadPlaceIcon(
	ctx context.Context,
	placeID uuid.UUID,
	key string,
) error {
	if err := validateTempPlaceIconKey(placeID, key); err != nil {
		return err
	}

	if err := s.s3.DeleteObject(ctx, key); err != nil {
		return fmt.Errorf("deleting temp place icon object: %w", err)
	}

	return nil
}

func (s *Uploader) DeletePlaceIcon(
	ctx context.Context,
	placeID uuid.UUID,
	key string,
) error {
	if err := validateFinalPlaceIconKey(placeID, key); err != nil {
		return err
	}

	if err := s.s3.DeleteObject(ctx, key); err != nil {
		return fmt.Errorf("deleting place icon object: %w", err)
	}

	return nil
}

func createTempPlaceIconKey(placeID uuid.UUID) string {
	return fmt.Sprintf("place/icon/%s/temp/%s", placeID, uuid.New().String())
}

var tempPlaceIconKeyRe = regexp.MustCompile(
	`^place/icon/([0-9a-fA-F-]{36})/temp/([0-9a-fA-F-]{36})$`,
)

func validateTempPlaceIconKey(placeID uuid.UUID, key string) error {
	if key == "" {
		return errx.ErrorPlaceIconIsInvalid.Raise(fmt.Errorf("empty key"))
	}

	matches := tempPlaceIconKeyRe.FindStringSubmatch(key)
	if matches == nil {
		return errx.ErrorPlaceIconIsInvalid.Raise(fmt.Errorf("key %s does not match temp place icon key pattern", key))
	}

	if matches[1] != placeID.String() {
		return errx.ErrorPlaceIconIsInvalid.Raise(fmt.Errorf("key %s does not belong to place %s", key, placeID))
	}

	return nil
}

func createPlaceIconKey(placeID uuid.UUID) string {
	return fmt.Sprintf("place/icon/%s/%s", placeID, uuid.New().String())
}

var finalPlaceIconKeyRe = regexp.MustCompile(
	`^place/icon/([0-9a-fA-F-]{36})/([0-9a-fA-F-]{36})$`,
)

func validateFinalPlaceIconKey(placeID uuid.UUID, key string) error {
	if key == "" {
		return errx.ErrorPlaceIconIsInvalid.Raise(fmt.Errorf("empty key"))
	}

	matches := finalPlaceIconKeyRe.FindStringSubmatch(key)
	if matches == nil {
		return errx.ErrorPlaceIconIsInvalid.Raise(fmt.Errorf("key %s does not match final place icon key pattern", key))
	}

	if matches[1] != placeID.String() {
		return errx.ErrorPlaceIconIsInvalid.Raise(fmt.Errorf("key %s does not belong to place %s", key, placeID))
	}

	return nil
}
