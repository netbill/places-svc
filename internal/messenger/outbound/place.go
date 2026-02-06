package outbound

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/evebox/header"
	"github.com/netbill/places-svc/internal/core/models"
	"github.com/netbill/places-svc/internal/messenger/contracts"
	"github.com/segmentio/kafka-go"
)

func (p *Producer) PublishCreatePlace(ctx context.Context, place models.Place) error {
	payload, err := json.Marshal(contracts.PlaceCreatedPayload{
		PlaceID:        place.ID,
		ClassID:        place.ClassID,
		OrganizationID: place.OrganizationID,

		Status:   place.Status,
		Verified: place.Verified,
		Point:    place.Point,
		Address:  place.Address,
		Name:     place.Name,

		Description: place.Description,
		Icon:        place.Icon,
		Banner:      place.Banner,
		Website:     place.Website,
		Phone:       place.Phone,

		CreatedAt: place.CreatedAt,
		UpdatedAt: place.UpdatedAt,
	})
	if err != nil {
		return err
	}

	_, err = p.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.PlacesTopicV1,
			Key:   []byte(place.ID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.PlaceCreatedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.PlaceSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)

	return err
}

func (p *Producer) PublishUpdatePlace(ctx context.Context, place models.Place) error {
	payload, err := json.Marshal(contracts.PlaceUpdatedPayload{
		PlaceID:        place.ID,
		ClassID:        place.ClassID,
		OrganizationID: place.OrganizationID,

		Status:   place.Status,
		Verified: place.Verified,
		Address:  place.Address,
		Name:     place.Name,

		Description: place.Description,
		Icon:        place.Icon,
		Banner:      place.Banner,
		Website:     place.Website,
		Phone:       place.Phone,

		UpdatedAt: place.UpdatedAt,
	})
	if err != nil {
		return err
	}

	_, err = p.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.PlacesTopicV1,
			Key:   []byte(place.ID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.PlaceUpdatedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.PlaceSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)

	return err
}

func (p *Producer) PublishUpdatePlaceStatus(ctx context.Context, place models.Place) error {
	payload, err := json.Marshal(contracts.PlaceStatusUpdatedPayload{
		PlaceID:   place.ID,
		Status:    place.Status,
		UpdatedAt: place.UpdatedAt,
	})
	if err != nil {
		return err
	}

	_, err = p.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.PlacesTopicV1,
			Key:   []byte(place.ID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.PlaceStatusUpdatedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.PlaceSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)

	return err
}

func (p *Producer) PublishUpdatePlaceVerified(ctx context.Context, place models.Place) error {
	payload, err := json.Marshal(contracts.PlaceVerifiedUpdatedPayload{
		PlaceID:   place.ID,
		Verified:  place.Verified,
		UpdatedAt: place.UpdatedAt,
	})
	if err != nil {
		return err
	}

	_, err = p.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.PlacesTopicV1,
			Key:   []byte(place.ID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.PlaceVerifiedUpdatedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.PlaceSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)

	return err
}

func (p *Producer) PublishUpdatePlaceClassID(ctx context.Context, place models.Place) error {
	payload, err := json.Marshal(contracts.PlaceClassIDUpdatedPayload{
		PlaceID:   place.ID,
		ClassID:   place.ClassID,
		UpdatedAt: place.UpdatedAt,
	})
	if err != nil {
		return err
	}

	_, err = p.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.PlacesTopicV1,
			Key:   []byte(place.ID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.PlaceClassIDUpdatedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.PlaceSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)

	return err
}

func (p *Producer) PublishDeletePlace(ctx context.Context, placeID uuid.UUID) error {
	payload, err := json.Marshal(contracts.PlaceDeletedPayload{
		PlaceID:   placeID,
		DeletedAt: time.Now().UTC(),
	})
	if err != nil {
		return err
	}

	_, err = p.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.PlacesTopicV1,
			Key:   []byte(placeID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.PlaceDeletedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.PlaceSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)

	return err
}
