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
	CreateOrganization(ctx context.Context, params models.Organization) error
	UpdateOrganizationStatus(ctx context.Context, ID uuid.UUID, status string) (models.Organization, error)
	DeleteOrganization(ctx context.Context, ID uuid.UUID) error

	CreateMember(ctx context.Context, member models.Member) error
	DeleteMember(ctx context.Context, memberID uuid.UUID) error

	RemoveMemberRole(ctx context.Context, memberID, roleID uuid.UUID) error
	AddMemberRole(ctx context.Context, memberID, roleID uuid.UUID) error

	CreateRole(ctx context.Context, params models.Role) error
	UpdateRolesRanks(
		ctx context.Context,
		organizationID uuid.UUID,
		order map[uuid.UUID]uint,
	) error
	UpdateRolePermissions(
		ctx context.Context,
		roleID uuid.UUID,
		permissions map[string]bool,
	) error
	DeleteRole(ctx context.Context, roleID uuid.UUID) error

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
