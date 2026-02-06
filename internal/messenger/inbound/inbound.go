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
	log  *logium.Logger
	core core
}

type core struct {
	organization organizationSvc
}

func New(log *logium.Logger, org organizationSvc) *Inbound {
	return &Inbound{
		log: log,
		core: core{
			organization: org,
		},
	}
}

type organizationSvc interface {
	CreateOrganization(ctx context.Context, organization models.Organization) (models.Organization, error)
	UpdateOrganization(
		ctx context.Context,
		organizationID uuid.UUID,
		params organization.UpdateParams,
	) (models.Organization, error)
	ActivateOrganization(
		ctx context.Context,
		orgID uuid.UUID,
		updatedAt time.Time,
	) (models.Organization, error)
	DeactivateOrganization(
		ctx context.Context,
		orgID uuid.UUID,
		deactivatedAt time.Time,
	) (models.Organization, error)
	DeleteOrganization(ctx context.Context, organizationID uuid.UUID) error

	CreateOrgMember(ctx context.Context, member models.OrgMember) (models.OrgMember, error)
	DeleteOrgMember(ctx context.Context, ID uuid.UUID) error

	AddOrgMemberRole(ctx context.Context, memberID, roleID uuid.UUID, addedAT time.Time) error
	RemoveOrgMemberRole(ctx context.Context, memberID, roleID uuid.UUID) error

	CreateOrgRole(ctx context.Context, role models.OrgRole) (models.OrgRole, error)
	DeleteOrgRole(ctx context.Context, ID uuid.UUID) error

	UpdateOrgRolePermissions(
		ctx context.Context,
		roleID uuid.UUID,
		permissions ...string,
	) error

	UpdateOrgRolesRanks(
		ctx context.Context,
		organizationID uuid.UUID,
		order map[uuid.UUID]uint,
		updatedAt time.Time,
	) error
}
