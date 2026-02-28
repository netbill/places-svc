package handler

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/netbill/eventbox"
	"github.com/netbill/evtypes"
	"github.com/netbill/places-svc/internal/core/errx"
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

	log := h.log.WithInboxEvent(event)

	err := h.modules.Org.CreateOrgMember(ctx, organization.CreateMemberParams{
		ID:             payload.MemberID,
		AccountID:      payload.AccountID,
		OrganizationID: payload.OrganizationID,
		Head:           payload.Head,
		Label:          payload.Label,
		Position:       payload.Position,
		CreatedAt:      payload.CreatedAt,
	})
	switch {
	case errors.Is(err, errx.ErrorOrgMemberDeleted):
		log.Debug("received org member created event for already deleted org member")
		return nil
	case errors.Is(err, errx.ErrorOrganizationDeleted):
		log.Debug("received org member created event for already deleted organization")
		return nil
	case errors.Is(err, errx.ErrorOrgMemberAlreadyExists):
		log.Debug("received org member created event for already existing org member")
		return nil
	case err != nil:
		return err
	default:
		log.Debug("org member created successfully")
		return nil
	}
}

func (h *Handler) OrgMemberUpdated(
	ctx context.Context,
	event eventbox.InboxEvent,
) error {
	var payload evtypes.OrgMemberUpdatedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return err
	}

	err := h.modules.Org.UpdateOrgMember(ctx, payload.MemberID, organization.UpdateMemberParams{
		Label:     payload.Label,
		Position:  payload.Position,
		Version:   payload.Version,
		UpdatedAt: payload.UpdatedAt,
	})
	switch {
	case errors.Is(err, errx.ErrorOrgMemberDeleted):
		h.log.WithInboxEvent(event).Debug("received org member updated event for already deleted org member")
		return nil
	case err != nil:
		return err
	default:
		h.log.WithInboxEvent(event).Debug("org member updated successfully")
		return nil
	}
}

func (h *Handler) OrgMemberDeleted(
	ctx context.Context,
	event eventbox.InboxEvent,
) error {
	var payload evtypes.OrgMemberDeletedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return err
	}

	err := h.modules.Org.DeleteOrgMember(ctx, payload.MemberID)
	switch {
	case errors.Is(err, errx.ErrorOrgMemberDeleted):
		h.log.WithInboxEvent(event).Debug("received org member deleted event for already deleted org member")
		return nil
	case err != nil:
		return err
	default:
		h.log.WithInboxEvent(event).Debug("org member deleted successfully")
		return nil
	}
}
