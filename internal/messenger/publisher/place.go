package publisher

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/eventbox"
	"github.com/netbill/evtypes"
	"github.com/netbill/places-svc/internal/core/models"
)

func (p *Publisher) PublishCreatePlace(
	ctx context.Context,
	place models.Place,
) error {
	payload, err := json.Marshal(evtypes.PlaceCreatedPayload{
		PlaceID:        place.ID,
		ClassID:        place.ClassID,
		OrganizationID: place.OrganizationID,

		Status:   place.Status,
		Verified: place.Verified,
		Point:    place.Point,
		Address:  place.Address,
		Name:     place.Name,

		Description: place.Description,
		IconKey:     place.IconKey,
		BannerKey:   place.BannerKey,
		Website:     place.Website,
		Phone:       place.Phone,

		Version:   place.Version,
		CreatedAt: place.CreatedAt,
	})
	if err != nil {
		return err
	}

	_, err = p.outbox.WriteOutboxEvent(ctx, eventbox.Message{
		ID:       uuid.New(),
		Type:     evtypes.PlaceCreatedEvent,
		Version:  place.Version,
		Topic:    evtypes.PlacesTopicV1,
		Key:      place.ID.String(),
		Producer: p.identity,
		Payload:  payload,
	})

	return err
}

func (p *Publisher) PublishUpdatePlace(
	ctx context.Context,
	place models.Place,
) error {
	payload, err := json.Marshal(evtypes.PlaceUpdatedPayload{
		PlaceID:        place.ID,
		ClassID:        place.ClassID,
		OrganizationID: place.OrganizationID,

		Status:   place.Status,
		Verified: place.Verified,
		Address:  place.Address,
		Name:     place.Name,

		Description: place.Description,
		IconKey:     place.IconKey,
		BannerKey:   place.BannerKey,
		Website:     place.Website,
		Phone:       place.Phone,

		Version:   place.Version,
		UpdatedAt: place.UpdatedAt,
	})
	if err != nil {
		return err
	}

	_, err = p.outbox.WriteOutboxEvent(ctx, eventbox.Message{
		ID:       uuid.New(),
		Type:     evtypes.PlaceUpdatedEvent,
		Version:  place.Version,
		Topic:    evtypes.PlacesTopicV1,
		Key:      place.ID.String(),
		Producer: p.identity,
		Payload:  payload,
	})

	return err
}

func (p *Publisher) PublishDeletePlace(
	ctx context.Context,
	placeID uuid.UUID,
) error {
	payload, err := json.Marshal(evtypes.PlaceDeletedPayload{
		PlaceID:   placeID,
		DeletedAt: time.Now().UTC(),
	})
	if err != nil {
		return err
	}

	_, err = p.outbox.WriteOutboxEvent(ctx, eventbox.Message{
		ID:       uuid.New(),
		Type:     evtypes.PlaceDeletedEvent,
		Version:  1,
		Topic:    evtypes.PlacesTopicV1,
		Key:      placeID.String(),
		Producer: p.identity,
		Payload:  payload,
	})

	return err
}
