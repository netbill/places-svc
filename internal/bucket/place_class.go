package bucket

import (
	"context"
	"fmt"
	"io"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/models"
)

func CreateTempPlaceClassIconKey(placeClassID, sessionID uuid.UUID) string {
	return fmt.Sprintf("place_class/icon/%s/temp/%s", placeClassID, sessionID)
}

func CreatePlaceClassIconKey(placeClassID uuid.UUID) string {
	return fmt.Sprintf("place_class/icon/%s", placeClassID)
}

func (b Bucket) GeneratePreloadLinkForPlaceClassMedia(
	ctx context.Context,
	placeClassID, sessionID uuid.UUID,
) (models.PlaceClassUploadMediaLinks, error) {
	uploadIconURL, getIconURL, err := b.s3.PresignPut(
		ctx,
		CreateTempPlaceClassIconKey(placeClassID, sessionID),
		b.tokensTTL.Place,
	)
	if err != nil {
		return models.PlaceClassUploadMediaLinks{}, fmt.Errorf("presigning put for place class icon: %w", err)
	}

	return models.PlaceClassUploadMediaLinks{
		IconUploadURL: uploadIconURL,
		IconGetURL:    getIconURL,
	}, nil
}

func (b Bucket) AcceptUpdatePlaceClassMedia(
	ctx context.Context,
	placeClassID, sessionID uuid.UUID,
) (models.PlaceClassMedia, error) {
	iconFinalKey := CreatePlaceClassIconKey(placeClassID)
	iconTempKey := CreateTempPlaceClassIconKey(placeClassID, sessionID)

	icon, size, err := b.s3.GetObjectRange(ctx, iconTempKey, 2048)
	if err != nil {
		return models.PlaceClassMedia{}, fmt.Errorf("failed to get object range for place class icon: %w", err)
	}
	defer icon.Close()

	uploadIcon := false
	if size != 0 {
		uploadIcon = true
		if err = b.ValidatePlaceClassIconUpload(icon, size); err != nil {
			return models.PlaceClassMedia{}, err
		}
	}

	res := models.PlaceClassMedia{}

	if uploadIcon {
		link, err := b.s3.CopyObject(ctx, iconTempKey, iconFinalKey)
		if err != nil {
			return models.PlaceClassMedia{}, fmt.Errorf("failed to copy object for place class icon: %w", err)
		}
		res.Icon = &link
	}

	return res, nil
}

func (b Bucket) ValidatePlaceClassIconUpload(
	icon io.Reader,
	size int64,
) error {
	probe, err := io.ReadAll(icon)
	if err != nil {
		return fmt.Errorf("failed to read icon probe bytes: %w", err)
	}

	valid, err := b.placeClassIconValidator.ValidateImageSize(uint(size))
	if err != nil {
		return fmt.Errorf("failed to validate place class icon image size: %w", err)
	}
	if !valid {
		return errx.ErrorPlaceClassIconTooLarge.Raise(
			fmt.Errorf("uploaded place class icon size %d exceeds the maximum allowed size", size),
		)
	}

	valid, err = b.placeClassIconValidator.ValidateImageContentType(probe)
	if err != nil {
		return fmt.Errorf("failed to validate place class icon content type: %w", err)
	}
	if !valid {
		return errx.ErrorPlaceClassIconContentTypeNotAllowed.Raise(
			fmt.Errorf("uploaded place class icon content type is not allowed"),
		)
	}

	valid, err = b.placeClassIconValidator.ValidateImageFormat(probe)
	if err != nil {
		return fmt.Errorf("failed to validate place class icon image format: %w", err)
	}
	if !valid {
		return errx.ErrorPlaceClassIconContentFormatNotAllowed.Raise(
			fmt.Errorf("uploaded place class icon format is not allowed"),
		)
	}

	valid, err = b.placeClassIconValidator.ValidateImageResolution(probe)
	if err != nil {
		return fmt.Errorf("failed to validate place class icon image resolution: %w", err)
	}
	if !valid {
		return errx.ErrorPlaceClassIconResolutionNotAllowed.Raise(
			fmt.Errorf("uploaded place class icon has invalid image resolution"),
		)
	}

	return nil
}

func (b Bucket) CancelUpdatePlaceClassIcon(
	ctx context.Context,
	placeClassID, sessionID uuid.UUID,
) error {
	key := CreateTempPlaceClassIconKey(placeClassID, sessionID)

	if err := b.s3.DeleteObject(ctx, key); err != nil {
		return fmt.Errorf("deleting temp place class icon object: %w", err)
	}

	return nil
}

func (b Bucket) DeletePlaceClassIcon(
	ctx context.Context,
	placeClassID uuid.UUID,
) error {
	key := CreatePlaceClassIconKey(placeClassID)

	if err := b.s3.DeleteObject(ctx, key); err != nil {
		return fmt.Errorf("deleting place class icon object: %w", err)
	}

	return nil
}

func (b Bucket) CleanPlaceClassMediaSession(
	ctx context.Context,
	placeClassID, sessionID uuid.UUID,
) error {
	if err := b.s3.DeleteObject(ctx, CreateTempPlaceClassIconKey(placeClassID, sessionID)); err != nil {
		return fmt.Errorf("failed to delete temp object for place class icon: %w", err)
	}

	return nil
}
