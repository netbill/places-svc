package inbound

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/logium"
	"github.com/netbill/places-svc/internal/core/models"
	"github.com/netbill/places-svc/internal/core/modules/organization"
	"github.com/netbill/places-svc/internal/core/modules/profile"
)

type Inbound struct {
	log    logium.Logger
	domain domain
}

//type domain struct {
//	profileSvc
//	organizationSvc
//}

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
	UpdateProfile(ctx context.Context, accountID uuid.UUID, params profile.UpdateParams) error
	DeleteProfile(ctx context.Context, accountID uuid.UUID) error
}

type organizationSvc interface {
	CreateOrganization(ctx context.Context, params models.Organization) error
	UpdateOrganization(ctx context.Context, organizationID uuid.UUID, params organization.UpdateParams) error
	DeleteOrganization(ctx context.Context, organizationID uuid.UUID) error
	DeactivateOrganization(
		ctx context.Context,
		orgID uuid.UUID,
		updatedAt time.Time,
	) error

	CreateOrgMember(ctx context.Context, member models.OrgMember) error
	UpdateOrgMember(ctx context.Context, memberID uuid.UUID, params organization.UpdateMemberParams) (models.OrgMember, error)
	DeleteOrgMember(ctx context.Context, ID uuid.UUID) error

	AddOrgMemberRole(ctx context.Context, memberID, roleID uuid.UUID, addedAT time.Time) error
	RemoveOrgMemberRole(ctx context.Context, memberID, roleID uuid.UUID) error

	CreateOrgRole(ctx context.Context, role models.OrgRole) error
	DeleteOrgRole(ctx context.Context, ID uuid.UUID) error

	UpdateOrgRolePermissions(
		ctx context.Context,
		roleID uuid.UUID,
		permissions []models.OrgRolePermissionLink,
	) error

	UpdateOrgRolesRanks(
		ctx context.Context,
		organizationID uuid.UUID,
		order map[uuid.UUID]uint,
	) error
}
