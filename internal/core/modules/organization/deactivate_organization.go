package organization

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/models"
)

func (m *Module) DeactivateOrganization(
	ctx context.Context,
	orgID uuid.UUID,
	updatedAt time.Time,
) (org models.Organization, err error) {
	err = m.repo.Transaction(ctx, func(ctx context.Context) error {
		org, err = m.repo.UpdateOrgStatus(
			ctx,
			orgID,
			models.OrganizationStatusInactive,
			updatedAt,
		)
		if err != nil {
			return err
		}

		if err := m.repo.UpdatePlaceStatusForOrg(
			ctx,
			orgID,
			models.OrganizationStatusInactive,
		); err != nil {
			return err
		}

		return nil
	})

	return org, err
}
