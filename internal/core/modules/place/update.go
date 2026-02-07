package place

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/models"
)

func (m *Module) OpenUpdateSession(
	ctx context.Context,
	initiator models.AccountClaims,
	placeID uuid.UUID,
) (models.Place, models.UpdatePlaceMedia, error) {
	place, err := m.Get(ctx, placeID)
	if err != nil {
		return models.Place{}, models.UpdatePlaceMedia{}, err
	}

	if err = m.chekPermissionForUpdatePlace(ctx, initiator, place.ID); err != nil {
		return models.Place{}, models.UpdatePlaceMedia{}, err
	}

	uploadSessionID := uuid.New()

	links, err := m.bucket.GeneratePreloadLinkForPlaceMedia(ctx, place.ID, uploadSessionID)
	if err != nil {
		return models.Place{}, models.UpdatePlaceMedia{}, err
	}

	uploadToken, err := m.token.NewUploadPlaceMediaToken(
		initiator.GetAccountID(),
		placeID,
		uploadSessionID,
	)
	if err != nil {
		return models.Place{}, models.UpdatePlaceMedia{}, err
	}

	return place, models.UpdatePlaceMedia{
		Links: models.PlaceUploadMediaLinks{
			IconUploadURL:   links.IconUploadURL,
			IconGetURL:      links.IconGetURL,
			BannerUploadURL: links.BannerUploadURL,
			BannerGetURL:    links.BannerGetURL,
		},
		UploadSessionID: uploadSessionID,
		UploadToken:     uploadToken,
	}, nil
}

type UpdateParams struct {
	Address     string
	Name        string
	Description *string
	Website     *string
	Phone       *string

	Media UpdateMediaParams
}

type UpdateMediaParams struct {
	UploadSessionID uuid.UUID

	DeleteIcon   bool
	icon         *string
	DeleteBanner bool
	banner       *string
}

func (p UpdateParams) GetUpdatedIcon() *string {
	if p.Media.DeleteIcon {
		return nil
	}

	return p.Media.icon
}

func (p UpdateParams) GetUpdatedBanner() *string {
	if p.Media.DeleteBanner {
		return nil
	}

	return p.Media.banner
}

func (m *Module) ConfirmUpdateSession(
	ctx context.Context,
	initiator models.AccountClaims,
	placeID uuid.UUID,
	params UpdateParams,
) (place models.Place, err error) {
	place, err = m.repo.GetPlaceByID(ctx, placeID)
	if err != nil {
		return models.Place{}, err
	}

	if place.OrganizationID != nil {
		if err = m.chekPermissionForUpdatePlace(ctx, initiator, *place.OrganizationID); err != nil {
			return models.Place{}, err
		}
	}

	params.Media.icon = place.Icon
	params.Media.banner = place.Banner

	if params.Media.DeleteIcon {
		if err = m.bucket.DeletePlaceIcon(ctx, placeID); err != nil {
			return models.Place{}, err
		}
		params.Media.icon = nil
	}

	if params.Media.DeleteBanner {
		if err = m.bucket.DeletePlaceBanner(ctx, placeID); err != nil {
			return models.Place{}, err
		}
		params.Media.banner = nil
	}

	if !params.Media.DeleteIcon || !params.Media.DeleteBanner {
		links, err := m.bucket.AcceptUpdatePlaceMedia(
			ctx,
			placeID,
			params.Media.UploadSessionID,
		)
		if err != nil {
			return models.Place{}, err
		}

		params.Media.icon = links.Icon
		params.Media.banner = links.Banner
	}

	if err = m.bucket.CleanPlaceMediaSession(ctx, placeID, params.Media.UploadSessionID); err != nil {
		return models.Place{}, err
	}

	err = m.repo.Transaction(ctx, func(txCtx context.Context) error {
		place, err = m.repo.UpdatePlaceByID(txCtx, placeID, params)
		if err != nil {
			return err
		}

		if err = m.messenger.PublishUpdatePlace(txCtx, place); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return models.Place{}, err
	}

	return place, nil
}

func (m *Module) DeleteUpdateIconInSession(
	ctx context.Context,
	initiator models.AccountClaims,
	placeID, uploadSessionID uuid.UUID,
) error {
	if err := m.chekPermissionForUpdatePlace(ctx, initiator, placeID); err != nil {
		return err
	}

	return m.bucket.CancelUpdatePlaceIcon(
		ctx,
		placeID,
		uploadSessionID,
	)
}

func (m *Module) DeleteUpdateBannerInSession(
	ctx context.Context,
	initiator models.AccountClaims,
	placeID, uploadSessionID uuid.UUID,
) error {
	if err := m.chekPermissionForUpdatePlace(ctx, initiator, placeID); err != nil {
		return err
	}

	return m.bucket.CancelUpdatePlaceBanner(
		ctx,
		placeID,
		uploadSessionID,
	)
}

func (m *Module) CancelUpdate(
	ctx context.Context,
	initiator models.AccountClaims,
	placeID, uploadSessionID uuid.UUID,
) error {
	if err := m.chekPermissionForUpdatePlace(ctx, initiator, placeID); err != nil {
		return err
	}

	return m.bucket.CleanPlaceMediaSession(
		ctx,
		placeID,
		uploadSessionID,
	)
}
