package inbound

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/netbill/evebox/box/inbox"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/models"
	"github.com/netbill/places-svc/internal/messenger/contracts"
)

func (i Inbound) OrgMemberCreated(
	ctx context.Context,
	event inbox.Event,
) inbox.EventStatus {
	var payload contracts.OrgMemberCreatedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		i.log.Errorf("bad payload for %s, key %s, id: %s, error: %v", event.Type, event.Key, event.ID, err)
		return inbox.EventStatusFailed
	}

	if _, err := i.core.organization.CreateOrgMember(ctx, models.OrgMember{
		ID:             payload.MemberID,
		AccountID:      payload.AccountID,
		OrganizationID: payload.OrganizationID,
		Head:           payload.Head,
		CreatedAt:      payload.CreatedAt,
	}); err != nil {
		switch {
		case errors.Is(err, errx.ErrorInternal):
			i.log.Errorf("failed to handle member created, key %s, id: %s, error: %v", event.Key, event.ID, err)
			return inbox.EventStatusPending
		default:
			i.log.Errorf("failed to handle member created, key %s, id: %s, error: %v", event.Key, event.ID, err)
			return inbox.EventStatusFailed
		}
	}

	return inbox.EventStatusProcessed
}

func (i Inbound) OrgMemberDeleted(
	ctx context.Context,
	event inbox.Event,
) inbox.EventStatus {
	var payload contracts.OrgMemberDeletedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		i.log.Errorf("bad payload for %s, key %s, id: %s, error: %v", event.Type, event.Key, event.ID, err)
		return inbox.EventStatusFailed
	}

	if err := i.core.organization.DeleteOrgMember(ctx, payload.MemberID); err != nil {
		switch {
		case errors.Is(err, errx.ErrorInternal):
			i.log.Errorf(
				"failed to handle member deleted, key %s, id: %s, error: %v",
				event.Key, event.ID, err,
			)
			return inbox.EventStatusPending
		default:
			i.log.Errorf("failed to handle member deleted, key %s, id: %s, error: %v", event.Key, event.ID, err)
			return inbox.EventStatusFailed
		}
	}

	return inbox.EventStatusProcessed
}

func (i Inbound) OrgMemberAddedRole(
	ctx context.Context,
	event inbox.Event,
) inbox.EventStatus {
	var payload contracts.OrgMemberRoleAddedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		i.log.Errorf("bad payload for %s, key %s, id: %s, error: %v", event.Type, event.Key, event.ID, err)
		return inbox.EventStatusFailed
	}

	if err := i.core.organization.AddOrgMemberRole(ctx, payload.MemberID, payload.RoleID, payload.AddedAt); err != nil {
		switch {
		case errors.Is(err, errx.ErrorInternal):
			i.log.Errorf(
				"failed to handle member added role, key %s, id: %s, error: %v",
				event.Key, event.ID, err,
			)
			return inbox.EventStatusPending
		default:
			i.log.Errorf("failed to handle member added role, key %s, id: %s, error: %v", event.Key, event.ID, err)
			return inbox.EventStatusFailed
		}
	}

	return inbox.EventStatusProcessed
}

func (i Inbound) OrgMemberRemovedRole(
	ctx context.Context,
	event inbox.Event,
) inbox.EventStatus {
	var payload contracts.OrgMemberRoleRemovedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		i.log.Errorf("bad payload for %s, key %s, id: %s, error: %v", event.Type, event.Key, event.ID, err)
		return inbox.EventStatusFailed
	}

	if err := i.core.organization.RemoveOrgMemberRole(ctx, payload.MemberID, payload.RoleID); err != nil {
		switch {
		case errors.Is(err, errx.ErrorInternal):
			i.log.Errorf(
				"failed to handle member removed role, key %s, id: %s, error: %v",
				event.Key, event.ID, err,
			)
			return inbox.EventStatusPending
		default:
			i.log.Errorf("failed to handle member removed role, key %s, id: %s, error: %v", event.Key, event.ID, err)
			return inbox.EventStatusFailed
		}
	}

	return inbox.EventStatusProcessed
}
