package contracts

import (
	"time"

	"github.com/google/uuid"
)

const OrgMemberTopicV1 = "organization.member.v1"

const OrgMemberCreatedEvent = "member.created"

type OrgMemberCreatedPayload struct {
	MemberID       uuid.UUID `json:"member_id"`
	AccountID      uuid.UUID `json:"account_id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Head           bool      `json:"head"`
	Position       *string   `json:"position,omitempty"`
	Label          *string   `json:"label,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

const OrgMemberDeletedEvent = "member.deleted"

type OrgMemberDeletedPayload struct {
	MemberID  uuid.UUID `json:"member_id"`
	DeletedAt time.Time `json:"deleted_at"`
}
