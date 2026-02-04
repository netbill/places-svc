package organization

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/models"
)

func (m *Module) SuspendOrganization(
	ctx context.Context,
	orgID uuid.UUID,
	updatedAt time.Time,
) error {
	return m.repo.Transaction(ctx, func(ctx context.Context) error {
		if err := m.repo.UpdateOrgStatus(
			ctx,
			orgID,
			models.OrganizationStatusSuspended,
			updatedAt,
		); err != nil {
			return err
		}

		if err := m.repo.UpdatePlaceStatusForOrg(
			ctx,
			orgID,
			models.OrganizationStatusSuspended,
		); err != nil {
			return err
		}

		return nil
	})
}
