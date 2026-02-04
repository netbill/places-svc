package inbound

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/netbill/evebox/box/inbox"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/models"
	"github.com/netbill/places-svc/internal/core/modules/profile"
	"github.com/netbill/places-svc/internal/messenger/contracts"
)

func (i Inbound) ProfileCreated(
	ctx context.Context,
	event inbox.Event,
) inbox.EventStatus {
	var payload contracts.ProfileCreatedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		i.log.Errorf("bad payload for %s, key %s, id: %s, error: %v", event.Type, event.Key, event.ID, err)
		return inbox.EventStatusFailed
	}

	if err := i.domain.CreateProfile(ctx, models.Profile{
		AccountID: payload.AccountID,
		Username:  payload.Username,
		CreatedAt: payload.CreatedAt,
	}); err != nil {
		switch {
		case errors.Is(err, errx.ErrorInternal):
			i.log.Errorf(
				"failed to create profile due to internal error, key %s, id: %s, error: %v",
				event.Key, event.ID, err,
			)
			return inbox.EventStatusPending
		default:
			i.log.Errorf("failed to create profile, key %s, id: %s, error: %v", event.Key, event.ID, err)
			return inbox.EventStatusFailed
		}
	}

	return inbox.EventStatusProcessed
}

func (i Inbound) ProfileDeleted(
	ctx context.Context,
	event inbox.Event,
) inbox.EventStatus {
	var payload contracts.ProfileDeletedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		i.log.Errorf("bad payload for %s, key %s, id: %s, error: %v", event.Type, event.Key, event.ID, err)
		return inbox.EventStatusFailed
	}

	if err := i.domain.DeleteProfile(ctx, payload.AccountID); err != nil {
		switch {
		case errors.Is(err, errx.ErrorInternal):
			i.log.Errorf(
				"failed to delete profile due to internal error, key %s, id: %s, error: %v",
				event.Key, event.ID, err,
			)
			return inbox.EventStatusPending
		default:
			i.log.Errorf("failed to delete profile, key %s, id: %s, error: %v", event.Key, event.ID, err)
			return inbox.EventStatusFailed
		}
	}

	return inbox.EventStatusProcessed
}

func (i Inbound) ProfileUpdated(
	ctx context.Context,
	event inbox.Event,
) inbox.EventStatus {
	var payload contracts.ProfileUpdatedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		i.log.Errorf("bad payload for %s, key: %s, id: %s, error: %v", event.Type, event.Key, event.ID, err)
		return inbox.EventStatusFailed
	}

	err := i.domain.UpdateProfile(ctx, payload.AccountID, profile.UpdateParams{
		Username:  payload.Username,
		Official:  payload.Official,
		Pseudonym: payload.Pseudonym,
		Avatar:    payload.Avatar,
		UpdatedAt: payload.UpdatedAt,
	})
	if err != nil {
		switch {
		case errors.Is(err, errx.ErrorInternal):
			i.log.Errorf(
				"failed to upsert profile due to internal error, key: %s, id: %s, error: %v",
				event.Key, event.ID, err,
			)
			return inbox.EventStatusPending
		default:
			i.log.Errorf("failed to upsert profile, key: %s, id: %s, error: %v", event.Key, event.ID, err)
			return inbox.EventStatusFailed
		}
	}

	return inbox.EventStatusProcessed
}
