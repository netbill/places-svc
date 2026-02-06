package organization

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type UpdateOrgRolePermissionsParams struct {
	PlaceCreate UpdateOrgRolePermissionsDetails
	PlaceDelete UpdateOrgRolePermissionsDetails
	PlaceUpdate UpdateOrgRolePermissionsDetails
}

type UpdateOrgRolePermissionsDetails struct {
	Enable    bool
	CreatedAt time.Time
}

func (m *Module) UpdateOrgRolePermissions(
	ctx context.Context,
	roleID uuid.UUID,
	permissions UpdateOrgRolePermissionsParams,
) error {
	return m.repo.SetOrgRolePermissions(ctx, roleID, permissions)
}
