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

func (s *Uploader) CreatePlaceClassIconUploadMediaLinks(
	ctx context.Context,
	classID uuid.UUID,
) (models.UploadMediaLink, error) {
	key := createTempPlaceClassIconKey(classID)

	uploadURL, getURL, err := s.s3.PresignPut(ctx, key, s.config.LinkTTL)
	if err != nil {
		return models.UploadMediaLink{}, fmt.Errorf("presigning put for place class icon: %w", err)
	}

	return models.UploadMediaLink{
		Key:        key,
		PreloadUrl: getURL,
		UploadURL:  uploadURL,
	}, nil
}

func (s *Uploader) UpdatePlaceClassIcon(
	ctx context.Context,
	classID uuid.UUID,
	key string,
) (string, error) {
	err := validateTempPlaceClassIconKey(classID, key)
	if err != nil {
		return "", err
	}

	out, err := s.s3.GetObjectRange(ctx, key, 64*1024)
	switch {
	case errors.Is(err, awsx.ErrNotFound):
		return "", errx.ErrorPlaceClassIconIsInvalid.Raise(
			fmt.Errorf("place icon not found for key: %s", key),
		)
	case err != nil:
		return "", fmt.Errorf("get object range for place class icon: %w", err)
	}
	defer out.Body.Close()

	if err = s.config.PlaceClassIcon.Validate(out); err != nil {
		return "", errx.ErrorPlaceClassIconIsInvalid.Raise(
			fmt.Errorf("validating place icon content for key %s: %w", key, err),
		)
	}

	finalKey := createPlaceClassIconKey(classID)

	if err = s.s3.CopyObject(ctx, key, finalKey); err != nil {
		return "", fmt.Errorf("copying object for place class icon: %w", err)
	}

	return finalKey, nil
}

func (s *Uploader) DeleteUploadPlaceClassIcon(
	ctx context.Context,
	classID uuid.UUID,
	key string,
) error {
	if err := validateTempPlaceClassIconKey(classID, key); err != nil {
		return err
	}

	if err := s.s3.DeleteObject(ctx, key); err != nil {
		return fmt.Errorf("deleting temp place icon object: %w", err)
	}

	return nil
}

func (s *Uploader) DeletePlaceClassIcon(
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

func createTempPlaceClassIconKey(classID uuid.UUID) string {
	return fmt.Sprintf("place_class/icon/%s/temp/%s", classID, uuid.New().String())
}

var tempPlaceClassIconKeyRe = regexp.MustCompile(
	`^place_class/icon/([0-9a-fA-F-]{36})/temp/([0-9a-fA-F-]{36})$`,
)

func validateTempPlaceClassIconKey(classID uuid.UUID, key string) error {
	matches := tempPlaceClassIconKeyRe.FindStringSubmatch(key)
	if matches == nil {
		return errx.ErrorPlaceClassIconIsInvalid.Raise(
			fmt.Errorf("key %s does not match temp place class icon key pattern", key),
		)
	}

	if matches[1] != classID.String() {
		return errx.ErrorPlaceClassIconIsInvalid.Raise(
			fmt.Errorf("key %s does not belong to place class %s", key, classID),
		)
	}

	return nil
}

func createPlaceClassIconKey(classID uuid.UUID) string {
	return fmt.Sprintf("place_class/icon/%s/%s", classID, uuid.New().String())
}

var finalPlaceClassIconKeyRe = regexp.MustCompile(
	`^place_class/icon/([0-9a-fA-F-]{36})/([0-9a-fA-F-]{36})$`,
)

func validateFinalPlaceClassIconKey(classID uuid.UUID, key string) error {
	matches := finalPlaceClassIconKeyRe.FindStringSubmatch(key)
	if matches == nil {
		return errx.ErrorPlaceClassIconIsInvalid.Raise(
			fmt.Errorf("key %s does not match final place class icon key pattern", key),
		)
	}

	if matches[1] != classID.String() {
		return errx.ErrorPlaceClassIconIsInvalid.Raise(
			fmt.Errorf("key %s does not belong to place class %s", key, classID),
		)
	}

	return nil
}
