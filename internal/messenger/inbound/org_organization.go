package inbound

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/netbill/evebox/box/inbox"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/messenger/contracts"
)

func (i Inbound) OrganizationCreated(
	ctx context.Context,
	event inbox.Event,
) inbox.EventStatus {
	var payload contracts.OrganizationCreatedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		i.log.Errorf("bad payload for %s, key %s, id: %s, error: %v", event.Type, event.Key, event.ID, err)
		return inbox.EventStatusFailed
	}

	if err := i.domain.CreateOrganization(ctx, payload.Organization); err != nil {
		switch {
		case errors.Is(err, errx.ErrorInternal):
			i.log.Errorf(
				"failed to handle organization created due to internal error, key %s, id: %s, error: %v",
				event.Key, event.ID, err,
			)
			return inbox.EventStatusPending
		default:
			i.log.Errorf(
				"failed to handle organization created, key %s, id: %s, error: %v",
				event.Key, event.ID, err,
			)
			return inbox.EventStatusFailed
		}
	}

	return inbox.EventStatusProcessed
}

func (i Inbound) OrganizationDeleted(
	ctx context.Context,
	event inbox.Event,
) inbox.EventStatus {
	var payload contracts.OrganizationDeletedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		i.log.Errorf("bad payload for %s, key %s, id: %s, error: %v", event.Type, event.Key, event.ID, err)
		return inbox.EventStatusFailed
	}

	if err := i.domain.DeleteOrganization(ctx, payload.Organization.ID); err != nil {
		switch {
		case errors.Is(err, errx.ErrorInternal):
			i.log.Errorf(
				"failed to handle organization deleted due to internal error, key %s, id: %s, error: %v",
				event.Key, event.ID, err,
			)
			return inbox.EventStatusPending
		default:
			i.log.Errorf(
				"failed to handle organization deleted, key %s, id: %s, error: %v",
				event.Key, event.ID, err,
			)
			return inbox.EventStatusFailed
		}
	}

	return inbox.EventStatusProcessed
}

func (i Inbound) OrganizationActivated(
	ctx context.Context,
	event inbox.Event,
) inbox.EventStatus {
	var payload contracts.OrganizationActivatedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		i.log.Errorf("bad payload for %s, key %s, id: %s, error: %v", event.Type, event.Key, event.ID, err)
		return inbox.EventStatusFailed
	}

	if _, err := i.domain.UpdateOrganizationStatus(
		ctx,
		payload.Organization.ID,
		payload.Organization.Status,
	); err != nil {
		switch {
		case errors.Is(err, errx.ErrorInternal):
			i.log.Errorf(
				"failed to handle organization activated due to internal error, key %s, id: %s, error: %v",
				event.Key, event.ID, err,
			)
			return inbox.EventStatusPending
		default:
			i.log.Errorf(
				"failed to handle organization activated, key %s, id: %s, error: %v",
				event.Key, event.ID, err,
			)
			return inbox.EventStatusFailed
		}
	}

	return inbox.EventStatusProcessed
}

func (i Inbound) OrganizationDeactivated(
	ctx context.Context,
	event inbox.Event,
) inbox.EventStatus {
	var payload contracts.OrganizationDeactivatedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		i.log.Errorf("bad payload for %s, key %s, id: %s, error: %v", event.Type, event.Key, event.ID, err)
		return inbox.EventStatusFailed
	}

	if _, err := i.domain.UpdateOrganizationStatus(
		ctx,
		payload.Organization.ID,
		payload.Organization.Status,
	); err != nil {
		switch {
		case errors.Is(err, errx.ErrorInternal):
			i.log.Errorf(
				"failed to handle organization deactivated due to internal error, key %s, id: %s, error: %v",
				event.Key, event.ID, err,
			)
			return inbox.EventStatusPending
		default:
			i.log.Errorf(
				"failed to handle organization deactivated, key %s, id: %s, error: %v",
				event.Key, event.ID, err,
			)
			return inbox.EventStatusFailed
		}
	}

	return inbox.EventStatusProcessed
}

func (i Inbound) OrganizationSuspended(
	ctx context.Context,
	event inbox.Event,
) inbox.EventStatus {
	var payload contracts.OrganizationSuspendedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		i.log.Errorf("bad payload for %s, key %s, id: %s, error: %v", event.Type, event.Key, event.ID, err)
		return inbox.EventStatusFailed
	}

	if _, err := i.domain.UpdateOrganizationStatus(
		ctx,
		payload.Organization.ID,
		payload.Organization.Status,
	); err != nil {
		switch {
		case errors.Is(err, errx.ErrorInternal):
			i.log.Errorf(
				"failed to handle organization suspended due to internal error, key %s, id: %s, error: %v",
				event.Key, event.ID, err,
			)
			return inbox.EventStatusPending
		default:
			i.log.Errorf(
				"failed to handle organization suspended, key %s, id: %s, error: %v",
				event.Key, event.ID, err,
			)
			return inbox.EventStatusFailed
		}
	}

	return inbox.EventStatusProcessed
}
