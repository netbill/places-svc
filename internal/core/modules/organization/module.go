package organization

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/models"
)

type Module struct {
	repo repo
}

func New(repo repo) *Module {
	return &Module{
		repo: repo,
	}
}

type repo interface {
	CreateOrganization(ctx context.Context, params models.Organization) error
	UpdateOrganization(ctx context.Context, orgID uuid.UUID, params UpdateParams) error
	UpdateOrgStatus(
		ctx context.Context,
		orgID uuid.UUID,
		status string,
		updatedAt time.Time,
	) error
	DeleteOrganization(ctx context.Context, ID uuid.UUID) error

	CreateOrgMember(ctx context.Context, member models.OrgMember) error
	UpdateOrgMember(ctx context.Context, orgMemberID uuid.UUID, member UpdateMemberParams) (models.OrgMember, error)
	DeleteOrgMember(ctx context.Context, memberID uuid.UUID) error

	RemoveOrgMemberRole(ctx context.Context, memberID, roleID uuid.UUID) error
	AddOrgMemberRole(ctx context.Context, memberID, roleID uuid.UUID, addedAt time.Time) error

	CreateOrgRole(ctx context.Context, params models.OrgRole) error
	DeleteOrgRole(ctx context.Context, roleID uuid.UUID) error

	UpdateOrgRolesRanks(
		ctx context.Context,
		organizationID uuid.UUID,
		order map[uuid.UUID]uint,
	) error
	UpdateOrgRolePermissions(
		ctx context.Context,
		roleID uuid.UUID,
		permissions []models.OrgRolePermissionLink,
	) error

	UpdatePlaceStatusForOrg(
		ctx context.Context,
		organizationID uuid.UUID,
		status string,
	) error

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
