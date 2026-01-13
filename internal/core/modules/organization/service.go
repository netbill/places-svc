package organization

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/models"
)

type Service struct {
	repo repo
}

func New(repo repo) Service {
	return Service{
		repo: repo,
	}
}

type repo interface {
	UpsertOrganization(ctx context.Context, params models.Organization) error
	DeleteOrganization(ctx context.Context, ID uuid.UUID) error

	UpsertOrgMember(ctx context.Context, member models.Member) error
	DeleteOrgMember(ctx context.Context, memberID uuid.UUID) error

	RemoveOrgMemberRole(ctx context.Context, memberID, roleID uuid.UUID) error
	AddOrgMemberRole(ctx context.Context, memberID, roleID uuid.UUID) error

	UpsertOrgRole(ctx context.Context, params models.OrgRole) error
	DeleteOrgRole(ctx context.Context, roleID uuid.UUID) error

	UpdateOrgRolesRanks(
		ctx context.Context,
		organizationID uuid.UUID,
		order map[uuid.UUID]uint,
	) error
	UpdateOrgRolePermissions(
		ctx context.Context,
		roleID uuid.UUID,
		permissions map[string]bool,
	) error

	UpdatePlaceStatusForOrg(ctx context.Context, organizationID uuid.UUID, status string) error

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
