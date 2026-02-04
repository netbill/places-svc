package bucket

import (
	"context"
	"fmt"
	"io"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/models"
)

func CreateTempPlaceBannerKey(placeID, sessionID uuid.UUID) string {
	return fmt.Sprintf("place/banner/%s/temp/%s", placeID, sessionID)
}

func CreatePlaceBannerKey(placeID uuid.UUID) string {
	return fmt.Sprintf("place/banner/%s", placeID)
}

func CreateTempPlaceIconKey(placeID, sessionID uuid.UUID) string {
	return fmt.Sprintf("place/icon/%s/temp/%s", placeID, sessionID)
}

func CreatePlaceIconKey(placeID uuid.UUID) string {
	return fmt.Sprintf("place/icon/%s", placeID)
}

func (b Bucket) GeneratePreloadLinkForPlaceMedia(
	ctx context.Context,
	placeID, sessionID uuid.UUID,
) (models.PlaceUploadMediaLinks, error) {
	uploadBannerURL, getBannerURL, err := b.s3.PresignPut(
		ctx,
		CreateTempPlaceBannerKey(placeID, sessionID),
		b.tokensTTL.Org,
	)
	if err != nil {
		return models.PlaceUploadMediaLinks{}, fmt.Errorf("presigning put for place banner: %w", err)
	}

	uploadIconURL, getIconURL, err := b.s3.PresignPut(
		ctx,
		CreateTempPlaceIconKey(placeID, sessionID),
		b.tokensTTL.Org,
	)
	if err != nil {
		return models.PlaceUploadMediaLinks{}, fmt.Errorf("presigning put for place icon: %w", err)
	}

	return models.PlaceUploadMediaLinks{
		BannerUploadURL: uploadBannerURL,
		BannerGetURL:    getBannerURL,
		IconUploadURL:   uploadIconURL,
		IconGetURL:      getIconURL,
	}, nil
}

func (b Bucket) AcceptUpdatePlaceMedia(
	ctx context.Context,
	placeID, sessionID uuid.UUID,
) (models.PlaceMedia, error) {
	IconFinalKey := CreatePlaceIconKey(placeID)
	IconTempKey := CreateTempPlaceIconKey(placeID, sessionID)

	icon, size, err := b.s3.GetObjectRange(ctx, IconTempKey, 2048)
	if err != nil {
		return models.PlaceMedia{}, fmt.Errorf("failed to get object range for place icon: %w", err)
	}
	defer icon.Close()

	uploadIcon := false

	if size != 0 {
		uploadIcon = true
		if err = b.ValidatePlaceIconUpload(icon, size); err != nil {
			return models.PlaceMedia{}, err
		}
	}

	BannerFinalKey := CreatePlaceBannerKey(placeID)
	BannerTempKey := CreateTempPlaceBannerKey(placeID, sessionID)

	banner, size, err := b.s3.GetObjectRange(ctx, BannerTempKey, 2048)
	if err != nil {
		return models.PlaceMedia{}, fmt.Errorf("failed to get object range for place banner: %w", err)
	}
	defer banner.Close()

	uploadBanner := false

	if size != 0 {
		uploadBanner = true
		if err = b.ValidatePlaceBannerUpload(banner, size); err != nil {
			return models.PlaceMedia{}, err
		}
	}

	res := models.PlaceMedia{}

	if uploadIcon {
		link, err := b.s3.CopyObject(ctx, IconTempKey, IconFinalKey)
		if err != nil {
			return models.PlaceMedia{}, fmt.Errorf("failed to copy object for place icon: %w", err)
		}
		res.Icon = &link
	}

	if uploadBanner {
		link, err := b.s3.CopyObject(ctx, BannerTempKey, BannerFinalKey)
		if err != nil {
			return models.PlaceMedia{}, fmt.Errorf("failed to copy object for place banner: %w", err)
		}
		res.Icon = &link
	}

	return res, nil
}

func (b Bucket) ValidatePlaceIconUpload(
	icon io.Reader,
	size int64,
) error {
	probe, err := io.ReadAll(icon)
	if err != nil {
		return fmt.Errorf("failed to read icon probe bytes: %w", err)
	}

	valid, err := b.OrgIconValidator.ValidateImageSize(uint(size))
	if err != nil {
		return fmt.Errorf("failed to validate place icon image size: %w", err)
	}
	if !valid {
		return errx.ErrorPlaceIconTooLarge.Raise(
			fmt.Errorf("uploaded place icon size %d exceeds the maximum allowed size", size),
		)
	}

	valid, err = b.OrgIconValidator.ValidateImageContentType(probe)
	if err != nil {
		return fmt.Errorf("failed to validate place icon content type: %w", err)
	}
	if !valid {
		return errx.ErrorPlaceIconContentTypeNotAllowed.Raise(
			fmt.Errorf("uploaded place icon content type is not allowed"),
		)
	}

	valid, err = b.OrgIconValidator.ValidateImageFormat(probe)
	if err != nil {
		return fmt.Errorf("failed to validate place icon image format: %w", err)
	}
	if !valid {
		return errx.ErrorPlaceIconContentFormatNotAllowed.Raise(
			fmt.Errorf("uploaded place icon format is not allowed"),
		)
	}

	valid, err = b.OrgIconValidator.ValidateImageResolution(probe)
	if err != nil {
		return fmt.Errorf("failed to validate place icon image resolution: %w", err)
	}
	if !valid {
		return errx.ErrorPlaceIconResolutionNotAllowed.Raise(
			fmt.Errorf("uploaded place icon has invalid image resolution"),
		)
	}

	return nil
}

func (b Bucket) ValidatePlaceBannerUpload(
	banner io.Reader,
	size int64,
) error {
	probe, err := io.ReadAll(banner)
	if err != nil {
		return fmt.Errorf("failed to read banner probe bytes: %w", err)
	}

	valid, err := b.OrgBannerValidator.ValidateImageSize(uint(size))
	if err != nil {
		return fmt.Errorf("failed to validate place banner image size: %w", err)
	}
	if !valid {
		return errx.ErrorPlaceBannerTooLarge.Raise(
			fmt.Errorf("uploaded place banner size %d exceeds the maximum allowed size", size),
		)
	}

	valid, err = b.OrgBannerValidator.ValidateImageContentType(probe)
	if err != nil {
		return fmt.Errorf("failed to validate place banner content type: %w", err)
	}
	if !valid {
		return errx.ErrorPlaceBannerContentTypeNotAllowed.Raise(
			fmt.Errorf("uploaded place banner content type is not allowed"),
		)
	}

	valid, err = b.OrgBannerValidator.ValidateImageFormat(probe)
	if err != nil {
		return fmt.Errorf("failed to validate place banner image format: %w", err)
	}
	if !valid {
		return errx.ErrorPlaceBannerContentFormatNotAllowed.Raise(
			fmt.Errorf("uploaded place banner format is not allowed"),
		)
	}

	valid, err = b.OrgBannerValidator.ValidateImageResolution(probe)
	if err != nil {
		return fmt.Errorf("failed to validate place banner image resolution: %w", err)
	}
	if !valid {
		return errx.ErrorPlaceBannerResolutionNotAllowed.Raise(
			fmt.Errorf("uploaded place banner has invalid image resolution"),
		)
	}

	return nil
}

func (b Bucket) CancelUpdatePlaceIcon(
	ctx context.Context,
	placeID, sessionID uuid.UUID,
) error {
	key := CreateTempPlaceIconKey(placeID, sessionID)

	if err := b.s3.DeleteObject(ctx, key); err != nil {
		return fmt.Errorf("deleting temp place icon object: %w", err)
	}

	return nil
}

func (b Bucket) CancelUpdatePlaceBanner(
	ctx context.Context,
	placeID, sessionID uuid.UUID,
) error {
	key := CreateTempPlaceBannerKey(placeID, sessionID)

	if err := b.s3.DeleteObject(ctx, key); err != nil {
		return fmt.Errorf("deleting temp place banner object: %w", err)
	}

	return nil
}

func (b Bucket) DeletePlaceIcon(
	ctx context.Context,
	placeID uuid.UUID,
) error {
	key := CreatePlaceIconKey(placeID)
	if err := b.s3.DeleteObject(ctx, key); err != nil {
		return fmt.Errorf("deleting place icon object: %w", err)
	}

	return nil
}

func (b Bucket) DeletePlaceBanner(
	ctx context.Context,
	placeID uuid.UUID,
) error {
	key := CreatePlaceBannerKey(placeID)
	if err := b.s3.DeleteObject(ctx, key); err != nil {
		return fmt.Errorf("deleting place banner object: %w", err)
	}

	return nil
}

func (b Bucket) CleanPlaceMediaSession(
	ctx context.Context,
	placeID, sessionID uuid.UUID,
) error {
	err := b.s3.DeleteObject(ctx, CreateTempPlaceIconKey(placeID, sessionID))
	if err != nil {
		return fmt.Errorf(
			"failed to delete temp object for place icon: %w", err,
		)
	}

	err = b.s3.DeleteObject(ctx, CreateTempPlaceBannerKey(placeID, sessionID))
	if err != nil {
		return fmt.Errorf(
			"failed to delete temp object for place banner: %w", err,
		)
	}

	return nil
}
