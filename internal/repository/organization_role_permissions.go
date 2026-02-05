package repository

import (
	"context"
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
	Insert(ctx context.Context, input OrganizationRolePermissionLinkRow) (OrganizationRolePermissionLinkRow, error)

	Get(ctx context.Context) (OrganizationRolePermissionLinkRow, error)
	Select(ctx context.Context) ([]OrganizationRolePermissionLinkRow, error)

	UpdateMany(ctx context.Context) (int64, error)
	UpdateOne(ctx context.Context) (OrganizationRolePermissionLinkRow, error)

	UpdateSourceCreatedAt(v time.Time) OrgRolePermissionLinksQ

	FilterByRoleID(roleID ...uuid.UUID) OrgRolePermissionLinksQ
	FilterByPermissionCode(code ...string) OrgRolePermissionLinksQ

	Delete(ctx context.Context) error
}
