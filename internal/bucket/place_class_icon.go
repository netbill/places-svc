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

func CreateTempPlaceClassIconKey(placeClassID uuid.UUID) string {
	return fmt.Sprintf("place_class/icon/%s/temp/%s", placeClassID, uuid.New().String())
}

func CreatePlaceClassIconKey(placeClassID uuid.UUID) string {
	return fmt.Sprintf("place_class/icon/%s/%s", placeClassID, uuid.New().String())
}

func (s *Storage) CreatePlaceClassIconUploadMediaLinks(
	ctx context.Context,
	classID uuid.UUID,
) (models.UploadMediaLink, error) {
	key := CreateTempPlaceClassIconKey(classID)

	uploadURL, getURL, err := s.s3.PresignPut(
		ctx,
		key,
		s.config.LinkTTL,
	)
	if err != nil {
		return models.UploadMediaLink{}, fmt.Errorf("presigning put for place class icon: %w", err)
	}

	return models.UploadMediaLink{
		Key:        key,
		PreloadUrl: getURL,
		UploadURL:  uploadURL,
	}, nil
}

func (s *Storage) ValidatePlaceClassIcon(
	ctx context.Context,
	classID uuid.UUID,
	key string,
) error {
	err := validateTempPlaceClassIconKey(classID, key)
	if err != nil {
		return err
	}

	out, err := s.s3.GetObjectRange(ctx, key, 64*1024)
	switch {
	case err != nil && err.Error() == "object not found":
		return errx.ErrorNoContentUploaded.Raise(
			fmt.Errorf("place class icon not found for key: %s", key),
		)
	case err != nil:
		return fmt.Errorf("get object range for place class icon: %w", err)
	}
	defer out.Body.Close()

	if err = s.config.PlaceClassIcon.Validate(out); err != nil {
		switch {
		case errors.Is(err, awsx.ErrorNoContentUploaded):
			return errx.ErrorNoContentUploaded.Raise(err)
		case errors.Is(err, awsx.ErrorSizeExceedsMax):
			return errx.ErrorPlaceClassIconContentIsExceedsMax.Raise(err)
		case errors.Is(err, awsx.ErrorResolutionIsInvalid):
			return errx.ErrorPlaceClassIconResolutionIsInvalid.Raise(err)
		case errors.Is(err, awsx.ErrorFormatNotAllowed):
			return errx.ErrorPlaceClassIconFormatIsNotAllowed.Raise(err)
		default:
			return fmt.Errorf("validate place class icon content: %w", err)
		}
	}

	return nil
}

func (s *Storage) DeleteUploadPlaceClassIcon(
	ctx context.Context,
	classID uuid.UUID,
	key string,
) error {
	if err := validateTempPlaceClassIconKey(classID, key); err != nil {
		return err
	}

	if err := s.s3.DeleteObject(ctx, key); err != nil {
		return fmt.Errorf("deleting temp place class icon object: %w", err)
	}

	return nil
}

func (s *Storage) DeletePlaceClassIcon(
	ctx context.Context,
	classID uuid.UUID,
	key string,
) error {
	if err := validateFinalPlaceClassIconKey(classID, key); err != nil {
		return err
	}

	if err := s.s3.DeleteObject(ctx, key); err != nil {
		return fmt.Errorf("deleting place class icon object: %w", err)
	}

	return nil
}

func (s *Storage) UpdatePlaceClassIcon(
	ctx context.Context,
	classID uuid.UUID,
	key string,
) (string, error) {
	if err := validateTempPlaceClassIconKey(classID, key); err != nil {
		return "", err
	}

	finalKey := CreatePlaceClassIconKey(classID)

	if err := s.s3.CopyObject(ctx, key, finalKey); err != nil {
		return "", fmt.Errorf("copying object for place class icon: %w", err)
	}

	return finalKey, nil
}

var (
	tempPlaceClassIconKeyRe = regexp.MustCompile(
		`^place_class/icon/([0-9a-fA-F-]{36})/temp/([0-9a-fA-F-]{36})$`,
	)

	finalPlaceClassIconKeyRe = regexp.MustCompile(
		`^place_class/icon/([0-9a-fA-F-]{36})/([0-9a-fA-F-]{36})$`,
	)
)

func validateTempPlaceClassIconKey(classID uuid.UUID, key string) error {
	if key == "" {
		return errx.ErrorPlaceClassIconKeyIsInvalid.Raise(fmt.Errorf("empty key"))
	}

	matches := tempPlaceClassIconKeyRe.FindStringSubmatch(key)
	if matches == nil {
		return errx.ErrorPlaceClassIconKeyIsInvalid.Raise(fmt.Errorf("key %s does not match temp place class icon key pattern", key))
	}

	if matches[1] != classID.String() {
		return errx.ErrorPlaceClassIconKeyIsInvalid.Raise(fmt.Errorf("key %s does not belong to place class %s", key, classID))
	}

	return nil
}

func validateFinalPlaceClassIconKey(classID uuid.UUID, key string) error {
	if key == "" {
		return errx.ErrorPlaceClassIconKeyIsInvalid.Raise(fmt.Errorf("empty key"))
	}

	matches := finalPlaceClassIconKeyRe.FindStringSubmatch(key)
	if matches == nil {
		return errx.ErrorPlaceClassIconKeyIsInvalid.Raise(fmt.Errorf("key %s does not match final place class icon key pattern", key))
	}

	if matches[1] != classID.String() {
		return errx.ErrorPlaceClassIconKeyIsInvalid.Raise(fmt.Errorf("key %s does not belong to place class %s", key, classID))
	}

	return nil
}
