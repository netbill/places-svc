package inbound

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/netbill/evebox/box/inbox"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/messenger/contracts"
)

func (i Inbound) MemberCreated(
	ctx context.Context,
	event inbox.Event,
) inbox.EventStatus {
	var payload contracts.MemberCreatedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		i.log.Errorf("bad payload for %s, key %s, id: %s, error: %v", event.Type, event.Key, event.ID, err)
		return inbox.EventStatusFailed
	}

	if err := i.domain.CreateMember(ctx, payload.Member); err != nil {
		switch {
		case errors.Is(err, errx.ErrorInternal):
			i.log.Errorf(
				"failed to handle member created, key %s, id: %s, error: %v",
				event.Key, event.ID, err,
			)
			return inbox.EventStatusPending
		default:
			i.log.Errorf("failed to handle member created, key %s, id: %s, error: %v", event.Key, event.ID, err)
			return inbox.EventStatusFailed
		}
	}

	return inbox.EventStatusProcessed
}

func (i Inbound) MemberDeleted(
	ctx context.Context,
	event inbox.Event,
) inbox.EventStatus {
	var payload contracts.MemberDeletedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		i.log.Errorf("bad payload for %s, key %s, id: %s, error: %v", event.Type, event.Key, event.ID, err)
		return inbox.EventStatusFailed
	}

	if err := i.domain.DeleteMember(ctx, payload.Member.ID); err != nil {
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

func (i Inbound) MemberAddedRole(
	ctx context.Context,
	event inbox.Event,
) inbox.EventStatus {
	var payload contracts.MemberRoleAddedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		i.log.Errorf("bad payload for %s, key %s, id: %s, error: %v", event.Type, event.Key, event.ID, err)
		return inbox.EventStatusFailed
	}

	if err := i.domain.AddMemberRole(ctx, payload.MemberID, payload.RoleID); err != nil {
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

func (i Inbound) MemberRemovedRole(
	ctx context.Context,
	event inbox.Event,
) inbox.EventStatus {
	var payload contracts.MemberRoleRemovedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		i.log.Errorf("bad payload for %s, key %s, id: %s, error: %v", event.Type, event.Key, event.ID, err)
		return inbox.EventStatusFailed
	}

	if err := i.domain.RemoveMemberRole(ctx, payload.MemberID, payload.RoleID); err != nil {
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
