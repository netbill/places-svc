package contracts

import "github.com/netbill/places-svc/internal/core/models"

const OrgMemberCreatedEvent = "member.created"

type OrgMemberCreatedPayload struct {
	Member models.Member `json:"member"`
}

const OrgMemberUpdatedEvent = "member.updated"

type OrgMemberUpdatedPayload struct {
	Member models.Member `json:"member"`
}

const OrgMemberDeletedEvent = "member.deleted"

type OrgMemberDeletedPayload struct {
	Member models.Member `json:"member"`
}
