package handler

import (
	"context"
	"encoding/json"

	"github.com/netbill/eventbox"
	"github.com/netbill/evtypes"
	"github.com/netbill/places-svc/internal/core/modules/organization"
)

func (h *Handler) OrgMemberCreated(
	ctx context.Context,
	event eventbox.InboxEvent,
) error {
	var payload evtypes.OrgMemberCreatedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return err
	}

	return h.modules.Org.CreateOrgMember(ctx, organization.CreateMemberParams{
		ID:             payload.MemberID,
		AccountID:      payload.AccountID,
		OrganizationID: payload.OrganizationID,
		Head:           payload.Head,
		Label:          payload.Label,
		Position:       payload.Position,
		CreatedAt:      payload.CreatedAt,
	})
}

func (h *Handler) OrgMemberUpdated(
	ctx context.Context,
	event eventbox.InboxEvent,
) error {
	var payload evtypes.OrgMemberUpdatedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return err
	}

	return h.modules.Org.UpdateOrgMember(ctx, payload.MemberID, organization.UpdateMemberParams{
		Label:     payload.Label,
		Position:  payload.Position,
		Version:   payload.Version,
		UpdatedAt: payload.UpdatedAt,
	})
}

func (h *Handler) OrgMemberDeleted(
	ctx context.Context,
	event eventbox.InboxEvent,
) error {
	var payload evtypes.OrgMemberDeletedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return err
	}

	return h.modules.Org.DeleteOrgMember(ctx, payload.MemberID)
}
