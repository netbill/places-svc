package organization

import (
	"context"
	"fmt"

	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/models"
)

func (s Service) UpdateOrganizationStatus(ctx context.Context, org models.Organization) error {
	return s.repo.Transaction(ctx, func(ctx context.Context) error {
		if err := s.repo.UpsertOrganization(ctx, org); err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to activate organization: %w", err))
		}

		statusForPlaces := ""
		switch org.Status {
		case models.OrganizationStatusActive:
			statusForPlaces = models.PlaceStatusInactive
		case models.OrganizationStatusInactive:
			statusForPlaces = models.PlaceStatusInactive
		case models.OrganizationStatusSuspended:
			statusForPlaces = models.PlaceStatusSuspended
		default:
			return errx.ErrorInternal.Raise(
				fmt.Errorf("unknown organization status: %s", org.Status))
		}

		if err := s.repo.UpdatePlaceStatusForOrg(ctx, org.ID, statusForPlaces); err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to update organization places status: %w", err))
		}

		return nil
	})
}
