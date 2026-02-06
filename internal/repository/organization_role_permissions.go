package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type OrganizationRolePermissionLinkRow struct {
	RoleID           uuid.UUID `db:"role_id"`
	PermissionCode   string    `db:"permission_code"`
	SourceCreatedAt  time.Time `db:"source_created_at"`
	ReplicaCreatedAt time.Time `db:"replica_created_at"`
}

type OrgRolePermissionLinksQ interface {
	New() OrgRolePermissionLinksQ
	Upsert(
		ctx context.Context,
		roleID uuid.UUID,
		data ...string,
	) ([]OrganizationRolePermissionLinkRow, error)

	Get(ctx context.Context) (OrganizationRolePermissionLinkRow, error)
	Select(ctx context.Context) ([]OrganizationRolePermissionLinkRow, error)
	Exist(ctx context.Context) (bool, error)

	FilterByRoleID(roleID ...uuid.UUID) OrgRolePermissionLinksQ
	FilterByPermissionCode(code ...string) OrgRolePermissionLinksQ

	Delete(ctx context.Context) error
}

func (r *Repository) SetOrgRolePermissions(
	ctx context.Context,
	roleID uuid.UUID,
	codes ...string,
) error {
	rows, err := r.OrgRolePermissionLinksQ.Upsert(ctx, roleID, codes...)
	if err != nil {
		return fmt.Errorf("set role permissions: %w", err)
	}

	enabled := make(map[string]struct{}, len(rows))
	for i := range rows {
		enabled[rows[i].PermissionCode] = struct{}{}
	}

	return nil
}

func (r *Repository) CheckMemberHavePermission(
	ctx context.Context,
	memberID uuid.UUID,
	permissionCode string,
) (bool, error) {
	exists, err := r.OrgRolePermissionLinksQ.New().
		FilterByRoleID(memberID).
		FilterByPermissionCode(permissionCode).
		Exist(ctx)
	if err != nil {
		return false, fmt.Errorf("check member have permission: %w", err)
	}

	return exists, nil
}
