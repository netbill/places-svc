package outbound

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/netbill/evebox/header"
	"github.com/netbill/places-svc/internal/core/models"
	"github.com/netbill/places-svc/internal/messenger/contracts"
	"github.com/segmentio/kafka-go"
)

func (p Producer) PublishPlaceClassCreated(ctx context.Context, class models.PlaceClass) error {
	payload, err := json.Marshal(contracts.PlaceClassCreatedPayload{
		PlaceClass: class,
	})
	if err != nil {
		return err
	}

	_, err = p.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.PlacesTopicV1,
			Key:   []byte(class.ID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.PlaceClassCreatedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.PlaceSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)

	return err
}

func (p Producer) PublishPlaceClassUpdated(ctx context.Context, class models.PlaceClass) error {
	payload, err := json.Marshal(contracts.PlaceClassUpdatedPayload{
		PlaceClass: class,
	})
	if err != nil {
		return err
	}

	_, err = p.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.PlacesTopicV1,
			Key:   []byte(class.ID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.PlaceClassUpdatedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.PlaceSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)

	return err
}

func (p Producer) PublishPlaceClassParentUpdated(ctx context.Context, class models.PlaceClass) error {
	payload, err := json.Marshal(contracts.PlaceClassParentUpdatedPayload{
		PlaceClass: class,
	})
	if err != nil {
		return err
	}

	_, err = p.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.PlacesTopicV1,
			Key:   []byte(class.ID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.PlaceClassParentUpdatedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.PlaceSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)

	return err
}

func (p Producer) PublishPlaceClassDeleted(ctx context.Context, classID uuid.UUID) error {
	payload, err := json.Marshal(contracts.PlaceClassDeletedPayload{
		PlaceClassID: classID,
	})
	if err != nil {
		return err
	}

	_, err = p.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.PlacesTopicV1,
			Key:   []byte(classID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.PlaceClassDeletedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.PlaceSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)

	return err
}

func (p Producer) PublishPlacesClassReplaced(ctx context.Context, oldClassID, newClassID uuid.UUID) error {
	payload, err := json.Marshal(contracts.PlacesClassReplacedPayload{
		OldClassID: oldClassID,
		NewClassID: newClassID,
	})
	if err != nil {
		return err
	}

	_, err = p.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.PlacesTopicV1,
			Key:   []byte(oldClassID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.PlacesClassReplacedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.PlaceSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)

	return err
}
