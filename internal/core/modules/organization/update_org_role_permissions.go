package organization

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/models"
)

func (m *Module) UpdateOrgRolePermissions(
	ctx context.Context,
	roleID uuid.UUID,
	permissions []models.OrgRolePermissionLink,
) error {
	return m.repo.UpdateOrgRolePermissions(ctx, roleID, permissions)
}
