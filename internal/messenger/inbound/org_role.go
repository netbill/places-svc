package inbound

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/netbill/evebox/box/inbox"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/messenger/contracts"
)

func (i Inbound) RoleCreated(
	ctx context.Context,
	event inbox.Event,
) inbox.EventStatus {
	var payload contracts.RoleCreatedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		i.log.Errorf("bad payload for %s, key %s, id: %s, error: %v", event.Type, event.Key, event.ID, err)
		return inbox.EventStatusFailed
	}

	if err := i.domain.CreateRole(ctx, payload.Role); err != nil {
		switch {
		case errors.Is(err, errx.ErrorInternal):
			i.log.Errorf(
				"failed to handle role created due to internal error, key %s, id: %s, error: %v",
				event.Key, event.ID, err,
			)
			return inbox.EventStatusPending
		default:
			i.log.Errorf(
				"failed to handle role created, key %s, id: %s, error: %v",
				event.Key, event.ID, err,
			)
			return inbox.EventStatusFailed
		}
	}

	return inbox.EventStatusProcessed
}

func (i Inbound) RoleDeleted(
	ctx context.Context,
	event inbox.Event,
) inbox.EventStatus {
	var payload contracts.RoleDeletedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		i.log.Errorf("bad payload for %s, key %s, id: %s, error: %v", event.Type, event.Key, event.ID, err)
		return inbox.EventStatusFailed
	}

	if err := i.domain.DeleteRole(ctx, payload.Role.ID); err != nil {
		switch {
		case errors.Is(err, errx.ErrorInternal):
			i.log.Errorf(
				"failed to handle role deleted due to internal error, key %s, id: %s, error: %v",
				event.Key, event.ID, err,
			)
			return inbox.EventStatusPending
		default:
			i.log.Errorf(
				"failed to handle role deleted, key %s, id: %s, error: %v",
				event.Key, event.ID, err,
			)
			return inbox.EventStatusFailed
		}
	}

	return inbox.EventStatusProcessed
}

func (i Inbound) RolePermissionsUpdated(
	ctx context.Context,
	event inbox.Event,
) inbox.EventStatus {
	var payload contracts.RolePermissionsUpdatedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		i.log.Errorf("bad payload for %s, key %s, id: %s, error: %v", event.Type, event.Key, event.ID, err)
		return inbox.EventStatusFailed
	}

	if err := i.domain.UpdateRolePermissions(ctx, payload.RoleID, payload.Permissions); err != nil {
		switch {
		case errors.Is(err, errx.ErrorInternal):
			i.log.Errorf(
				"failed to handle update role permissions, key %s, id: %s, error: %v",
				event.Key, event.ID, err,
			)
			return inbox.EventStatusPending
		default:
			i.log.Errorf(
				"failed to handle update role permissions, key %s, id: %s, error: %v",
				event.Key, event.ID, err,
			)
			return inbox.EventStatusFailed
		}
	}

	return inbox.EventStatusProcessed
}

func (i Inbound) RolesRanksUpdated(
	ctx context.Context,
	event inbox.Event,
) inbox.EventStatus {
	var payload contracts.RolesRanksUpdatedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		i.log.Errorf("bad payload for %s, key %s, id: %s, error: %v", event.Type, event.Key, event.ID, err)
		return inbox.EventStatusFailed
	}

	if err := i.domain.UpdateRolesRanks(ctx, payload.OrganizationID, payload.Ranks); err != nil {
		switch {
		case errors.Is(err, errx.ErrorInternal):
			i.log.Errorf(
				"failed to handle update roles ranks, key %s, id: %s, error: %v",
				event.Key, event.ID, err,
			)
			return inbox.EventStatusPending
		default:
			i.log.Errorf(
				"failed to handle update roles ranks, key %s, id: %s, error: %v",
				event.Key, event.ID, err,
			)
			return inbox.EventStatusFailed
		}
	}

	return inbox.EventStatusProcessed
}
