package organization

import (
	"context"

	"github.com/google/uuid"
)

func (m *Module) UpdateOrgRolePermissions(
	ctx context.Context,
	roleID uuid.UUID,
	permissions ...string,
) error {
	return m.repo.SetOrgRolePermissions(ctx, roleID, permissions...)
}
