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

func CreateTempPlaceIconKey(placeID uuid.UUID) string {
	return fmt.Sprintf("place/icon/%s/temp/%s", placeID, uuid.New().String())
}

func CreatePlaceIconKey(placeID uuid.UUID) string {
	return fmt.Sprintf("place/icon/%s/%s", placeID, uuid.New().String())
}

func (s *Storage) CreatePlaceIconUploadMediaLinks(
	ctx context.Context,
	placeID uuid.UUID,
) (models.UploadMediaLink, error) {
	key := CreateTempPlaceIconKey(placeID)

	uploadURL, getURL, err := s.s3.PresignPut(
		ctx,
		key,
		s.config.LinkTTL,
	)
	if err != nil {
		return models.UploadMediaLink{}, fmt.Errorf("presigning put for place place icon: %w", err)
	}

	return models.UploadMediaLink{
		Key:        key,
		PreloadUrl: getURL,
		UploadURL:  uploadURL,
	}, nil
}

func (s *Storage) ValidatePlaceIcon(
	ctx context.Context,
	placeID uuid.UUID,
	key string,
) error {
	err := validateTempPlaceIconKey(placeID, key)
	if err != nil {
		return err
	}

	out, err := s.s3.GetObjectRange(ctx, key, 64*1024)
	switch {
	case errors.Is(err, awsx.ErrNotFound):
		return errx.ErrorNoContentUploaded.Raise(
			fmt.Errorf("place place icon not found for key: %s", key),
		)
	case err != nil:
		return fmt.Errorf("get object range for place place icon: %w", err)
	}
	defer out.Body.Close()

	if err = s.config.PlaceIcon.Validate(out); err != nil {
		switch {
		case errors.Is(err, awsx.ErrorNoContentUploaded):
			return errx.ErrorNoContentUploaded.Raise(err)
		case errors.Is(err, awsx.ErrorSizeExceedsMax):
			return errx.ErrorPlaceIconContentIsExceedsMax.Raise(err)
		case errors.Is(err, awsx.ErrorResolutionIsInvalid):
			return errx.ErrorPlaceIconResolutionIsInvalid.Raise(err)
		case errors.Is(err, awsx.ErrorFormatNotAllowed):
			return errx.ErrorPlaceIconFormatIsNotAllowed.Raise(err)
		default:
			return fmt.Errorf("validate place icon content: %w", err)
		}
	}

	return nil
}

func (s *Storage) DeleteUploadPlaceIcon(
	ctx context.Context,
	placeID uuid.UUID,
	key string,
) error {
	if err := validateTempPlaceIconKey(placeID, key); err != nil {
		return err
	}

	if err := s.s3.DeleteObject(ctx, key); err != nil {
		return fmt.Errorf("deleting temp place place icon object: %w", err)
	}

	return nil
}

func (s *Storage) DeletePlaceIcon(
	ctx context.Context,
	placeID uuid.UUID,
	key string,
) error {
	if err := validateFinalPlaceIconKey(placeID, key); err != nil {
		return err
	}

	if err := s.s3.DeleteObject(ctx, key); err != nil {
		return fmt.Errorf("deleting place place icon object: %w", err)
	}

	return nil
}

func (s *Storage) UpdatePlaceIcon(
	ctx context.Context,
	placeID uuid.UUID,
	key string,
) (string, error) {
	if err := validateTempPlaceIconKey(placeID, key); err != nil {
		return "", err
	}

	finalKey := CreatePlaceIconKey(placeID)

	if err := s.s3.CopyObject(ctx, key, finalKey); err != nil {
		return "", fmt.Errorf("copying object for place icon: %w", err)
	}

	return finalKey, nil
}

var (
	tempPlaceIconKeyRe = regexp.MustCompile(
		`^place/icon/([0-9a-fA-F-]{36})/temp/([0-9a-fA-F-]{36})$`,
	)

	finalPlaceIconKeyRe = regexp.MustCompile(
		`^place/icon/([0-9a-fA-F-]{36})/([0-9a-fA-F-]{36})$`,
	)
)

func validateTempPlaceIconKey(placeID uuid.UUID, key string) error {
	if key == "" {
		return errx.ErrorPlaceIconKeyIsInvalid.Raise(fmt.Errorf("empty key"))
	}

	matches := tempPlaceIconKeyRe.FindStringSubmatch(key)
	if matches == nil {
		return errx.ErrorPlaceIconKeyIsInvalid.Raise(fmt.Errorf("key %s does not match temp place place icon key pattern", key))
	}

	if matches[1] != placeID.String() {
		return errx.ErrorPlaceIconKeyIsInvalid.Raise(fmt.Errorf("key %s does not belong to place place %s", key, placeID))
	}

	return nil
}

func validateFinalPlaceIconKey(placeID uuid.UUID, key string) error {
	if key == "" {
		return errx.ErrorPlaceIconKeyIsInvalid.Raise(fmt.Errorf("empty key"))
	}

	matches := finalPlaceIconKeyRe.FindStringSubmatch(key)
	if matches == nil {
		return errx.ErrorPlaceIconKeyIsInvalid.Raise(fmt.Errorf("key %s does not match final place place icon key pattern", key))
	}

	if matches[1] != placeID.String() {
		return errx.ErrorPlaceIconKeyIsInvalid.Raise(fmt.Errorf("key %s does not belong to place place %s", key, placeID))
	}

	return nil
}
