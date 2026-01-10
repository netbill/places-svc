package inbound

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/logium"
	"github.com/netbill/places-svc/internal/core/models"
	"github.com/netbill/places-svc/internal/core/modules/profile"
)

type Inbound struct {
	log    logium.Logger
	domain domain
}

type domain struct {
	profileSvc
	organizationSvc
}

func New(log logium.Logger, prof profileSvc, org organizationSvc) Inbound {
	return Inbound{
		log: log,
		domain: domain{
			profileSvc:      prof,
			organizationSvc: org,
		},
	}
}

type profileSvc interface {
	CreateProfile(ctx context.Context, profile models.Profile) error
	DeleteProfile(ctx context.Context, accountID uuid.UUID) error
	UpdateUsername(ctx context.Context, accountID uuid.UUID, username string) (models.Profile, error)
	UpdateProfile(ctx context.Context, ID uuid.UUID, params profile.UpdateProfileParams) (models.Profile, error)
}

type organizationSvc interface {
	CreateOrganization(ctx context.Context, params models.Organization) error
	DeleteOrganization(ctx context.Context, ID uuid.UUID) error
	UpdateOrganizationStatus(ctx context.Context, ID uuid.UUID, status string) (models.Organization, error)

	CreateMember(ctx context.Context, member models.Member) error
	DeleteMember(ctx context.Context, ID uuid.UUID) error

	AddMemberRole(ctx context.Context, memberID, roleID uuid.UUID) error
	RemoveMemberRole(ctx context.Context, memberID, roleID uuid.UUID) error

	CreateRole(ctx context.Context, role models.Role) error
	DeleteRole(ctx context.Context, ID uuid.UUID) error

	UpdateRolePermissions(
		ctx context.Context,
		roleID uuid.UUID,
		permissions map[string]bool,
	) error

	UpdateRolesRanks(
		ctx context.Context,
		organizationID uuid.UUID,
		order map[uuid.UUID]uint,
	) error
}
