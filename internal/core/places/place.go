package places

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/errx"
	"github.com/netbill/places-svc/internal/models"
)

type UpdateParams struct {
	ClassID *uuid.UUID
	Name    *string
	Address *string

	Description *string
	Website     *string
	Phone       *string

	IconKey   *string
	BannerKey *string
}

func ptrEqual[T comparable](a, b *T) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

func (p *UpdateParams) HasChanges(place models.Place) bool {
	return !ptrEqual(p.ClassID, &place.ClassID) ||
		!ptrEqual(p.Name, &place.Name) ||
		!ptrEqual(p.Address, &place.Address) ||
		!ptrEqual(p.Description, place.Description) ||
		!ptrEqual(p.Website, place.Website) ||
		!ptrEqual(p.Phone, place.Phone) ||
		!ptrEqual(p.IconKey, place.IconKey) ||
		!ptrEqual(p.BannerKey, place.BannerKey)
}

func (s *Service) Update(
	ctx context.Context,
	actor models.AccountActor,
	placeID uuid.UUID,
	params UpdateParams,
) (place models.Place, err error) {
	place, err = s.place.Get(ctx, placeID)
	if err != nil {
		return models.Place{}, err
	}

	_, err = s.org.AuthorizeOrgHead(ctx, actor, place.OrganizationID)
	if err != nil {
		return models.Place{}, err
	}

	_, err = s.org.ValidateOrg(ctx, place.OrganizationID)
	if err != nil {
		return models.Place{}, err
	}

	if !params.HasChanges(place) {
		return place, nil
	}

	if params.ClassID != nil {
		class, err := s.class.Get(ctx, *params.ClassID)
		if err != nil {
			return models.Place{}, err
		}
		if class.DeprecatedAt != nil {
			return models.Place{}, errx.ErrorPlaceClassIsDeprecated.Raise(
				fmt.Errorf("place class %s is deprecated", params.ClassID),
			)
		}
	}

	switch {
	case params.IconKey != nil && *params.IconKey == "" && place.IconKey != nil:
		if err := s.media.DeletePlaceIcon(ctx, placeID, *place.IconKey); err != nil {
			return models.Place{}, fmt.Errorf("failed to delete place icon: %w", err)
		}
		params.IconKey = nil
	case params.IconKey != nil:
		iconKey, err := s.media.UpdatePlaceIcon(ctx, placeID, *params.IconKey)
		if err != nil {
			return models.Place{}, fmt.Errorf("failed to validate place icon: %w", err)
		}
		params.IconKey = &iconKey
	}

	switch {
	case params.BannerKey != nil && *params.BannerKey == "" && place.BannerKey != nil:
		if err := s.media.DeletePlaceBanner(ctx, placeID, *place.BannerKey); err != nil {
			return models.Place{}, fmt.Errorf("failed to delete place banner: %w", err)
		}
		params.BannerKey = nil
	case params.BannerKey != nil:
		bannerKey, err := s.media.UpdatePlaceBanner(ctx, placeID, *params.BannerKey)
		if err != nil {
			return models.Place{}, fmt.Errorf("failed to validate place banner: %w", err)
		}
		params.BannerKey = &bannerKey
	}

	err = s.tx.Transaction(ctx, func(ctx context.Context) error {
		place, err = s.place.Update(ctx, placeID, params)
		if err != nil {
			return err
		}

		return s.messenger.PublishUpdatePlace(ctx, place)
	})
	if err != nil {
		return models.Place{}, err
	}

	return place, nil
}

func (s *Service) Activate(
	ctx context.Context,
	actor models.AccountActor,
	placeID uuid.UUID,
) (place models.Place, err error) {
	return s.updateStatus(ctx, actor, placeID, models.PlaceStatusActive)
}

func (s *Service) Deactivate(
	ctx context.Context,
	actor models.AccountActor,
	placeID uuid.UUID,
) (place models.Place, err error) {
	return s.updateStatus(ctx, actor, placeID, models.PlaceStatusInactive)
}

func (s *Service) updateStatus(
	ctx context.Context,
	actor models.AccountActor,
	placeID uuid.UUID,
	status string,
) (place models.Place, err error) {
	place, err = s.place.Get(ctx, placeID)
	if err != nil {
		return models.Place{}, err
	}

	_, err = s.org.AuthorizeOrgHead(ctx, actor, place.OrganizationID)
	if err != nil {
		return models.Place{}, err
	}

	_, err = s.org.ValidateOrg(ctx, place.OrganizationID)
	if err != nil {
		return models.Place{}, err
	}

	switch {
	case status == place.Status:
		return place, nil
	case status != models.PlaceStatusActive && status != models.PlaceStatusInactive:
		return models.Place{}, errx.ErrorPlaceStatusIsInvalid.Raise(
			fmt.Errorf("cannot set status to %s, suspended", status),
		)
	}

	err = s.tx.Transaction(ctx, func(ctx context.Context) error {
		place, err = s.place.UpdateStatus(ctx, placeID, status)
		if err != nil {
			return err
		}

		return s.messenger.PublishUpdatePlace(ctx, place)
	})
	if err != nil {
		return models.Place{}, err
	}

	return place, nil
}

func (s *Service) Verify(
	ctx context.Context,
	placeID uuid.UUID,
) (place models.Place, err error) {
	return s.updateVerified(ctx, placeID, true)
}

func (s *Service) Unverify(
	ctx context.Context,
	placeID uuid.UUID,
) (place models.Place, err error) {
	return s.updateVerified(ctx, placeID, false)
}

func (s *Service) updateVerified(
	ctx context.Context,
	placeID uuid.UUID,
	verified bool,
) (place models.Place, err error) {
	place, err = s.place.Get(ctx, placeID)
	if err != nil {
		return models.Place{}, err
	}

	if place.Verified == verified {
		return place, nil
	}

	err = s.tx.Transaction(ctx, func(ctx context.Context) error {
		place, err = s.place.UpdateVerified(ctx, placeID, verified)
		if err != nil {
			return err
		}

		return s.messenger.PublishUpdatePlace(ctx, place)
	})

	return place, nil
}
