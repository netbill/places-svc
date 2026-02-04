package organization

import (
	"context"

	"github.com/google/uuid"
)

func (m *Module) UpdateOrgRolesRanks(
	ctx context.Context,
	organizationID uuid.UUID,
	order map[uuid.UUID]uint,
) error {
	return m.repo.UpdateOrgRolesRanks(ctx, organizationID, order)
}
