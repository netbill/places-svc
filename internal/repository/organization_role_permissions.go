package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/models"
	"github.com/netbill/places-svc/internal/core/modules/organization"
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
		data map[string]time.Time,
	) ([]OrganizationRolePermissionLinkRow, error)

	Get(ctx context.Context) (OrganizationRolePermissionLinkRow, error)
	Select(ctx context.Context) ([]OrganizationRolePermissionLinkRow, error)
	Exist(ctx context.Context) (bool, error)

	UpdateMany(ctx context.Context) (int64, error)
	UpdateOne(ctx context.Context) (OrganizationRolePermissionLinkRow, error)

	UpdateSourceCreatedAt(v time.Time) OrgRolePermissionLinksQ

	FilterByRoleID(roleID ...uuid.UUID) OrgRolePermissionLinksQ
	FilterByPermissionCode(code ...string) OrgRolePermissionLinksQ

	Delete(ctx context.Context) error
}

func (r *Repository) SetOrgRolePermissions(
	ctx context.Context,
	roleID uuid.UUID,
	permissions organization.UpdateOrgRolePermissionsParams,
) error {
	codes := make(map[string]time.Time)

	if permissions.PlaceCreate.Enable {
		codes[models.RolePermissionPlaceCreate] = permissions.PlaceCreate.CreatedAt
	}
	if permissions.PlaceDelete.Enable {
		codes[models.RolePermissionPlaceDelete] = permissions.PlaceDelete.CreatedAt
	}
	if permissions.PlaceUpdate.Enable {
		codes[models.RolePermissionPlaceUpdate] = permissions.PlaceUpdate.CreatedAt
	}

	rows, err := r.OrgRolePermissionLinksQ.Upsert(ctx, roleID, codes)
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
