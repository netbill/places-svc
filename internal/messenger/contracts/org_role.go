package contracts

import (
	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/models"
)

const OrgRoleCreatedEvent = "role.created"

type OrgRoleCreatedPayload struct {
	Role models.OrgRole `json:"role"`
}

const OrgRoleUpdatedEvent = "role.updated"

type OrgRoleUpdatedPayload struct {
	Role models.OrgRole `json:"role"`
}

const OrgRoleDeletedEvent = "role.deleted"

type OrgRoleDeletedPayload struct {
	Role models.OrgRole `json:"role"`
}

const OrgRolesRanksUpdatedEvent = "roles.ranks.updated"

type OrgRolesRanksUpdatedPayload struct {
	OrganizationID uuid.UUID          `json:"organization_id"`
	Ranks          map[uuid.UUID]uint `json:"ranks"`
}

const OrgRolePermissionsUpdatedEvent = "role.permissions.updated"

type OrgRolePermissionsUpdatedPayload struct {
	RoleID      uuid.UUID       `json:"role_id"`
	Permissions map[string]bool `json:"permissions"`
}

const OrgMemberRoleAddedEvent = "member_role.added"

type OrgMemberRoleAddedPayload struct {
	MemberID uuid.UUID `json:"member_id"`
	RoleID   uuid.UUID `json:"role_id"`
}

const OrgMemberRoleRemovedEvent = "member_role.remove"

type OrgMemberRoleRemovedPayload struct {
	MemberID uuid.UUID `json:"member_id"`
	RoleID   uuid.UUID `json:"role_id"`
}
