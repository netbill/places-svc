package inbound

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/logium"
	"github.com/netbill/places-svc/internal/core/models"
	"github.com/netbill/places-svc/internal/core/modules/organization"
)

type Inbound struct {
	log    *logium.Logger
	domain domain
}

type domain struct {
	organizationSvc
}

func New(log *logium.Logger, org organizationSvc) Inbound {
	return Inbound{
		log: log,
		domain: domain{
			organizationSvc: org,
		},
	}
}

type organizationSvc interface {
	CreateOrganization(ctx context.Context, organization models.Organization) error
	UpdateOrganization(
		ctx context.Context,
		organizationID uuid.UUID,
		params organization.UpdateParams,
	) (models.Organization, error)
	ActivateOrganization(
		ctx context.Context,
		orgID uuid.UUID,
		updatedAt time.Time,
	) error
	DeactivateOrganization(
		ctx context.Context,
		orgID uuid.UUID,
		deactivatedAt time.Time,
	) error
	DeleteOrganization(ctx context.Context, organizationID uuid.UUID) error

	CreateOrgMember(ctx context.Context, member models.OrgMember) error
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
