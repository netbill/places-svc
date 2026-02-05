package organization

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/models"
)

func (m *Module) ActivateOrganization(
	ctx context.Context,
	orgID uuid.UUID,
	updatedAt time.Time,
) error {
	return m.repo.UpdateOrgStatus(
		ctx,
		orgID,
		models.OrganizationStatusInactive,
		updatedAt,
	)
}
