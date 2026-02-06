package organization

import (
	"context"
	"time"

	"github.com/google/uuid"
)

func (m *Module) UpdateOrgRolesRanks(
	ctx context.Context,
	organizationID uuid.UUID,
	order map[uuid.UUID]uint,
	updatedAt time.Time,
) error {
	return m.repo.UpdateOrgRolesRanks(ctx, organizationID, order, updatedAt)
}
