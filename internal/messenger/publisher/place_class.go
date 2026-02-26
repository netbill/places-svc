package publisher

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/netbill/eventbox"
	"github.com/netbill/evtypes"
	"github.com/netbill/places-svc/internal/core/models"
)

func (p *Publisher) PublishPlaceClassCreated(
	ctx context.Context,
	class models.PlaceClass,
) error {
	payload, err := json.Marshal(evtypes.PlaceClassCreatedPayload{
		PlaceClassID: class.ID,
		ParentID:     class.ParentID,
		Name:         class.Name,
		Description:  class.Description,
		IconKey:      class.IconKey,
		Version:      class.Version,
		CreatedAt:    class.CreatedAt,
	})
	if err != nil {
		return err
	}

	_, err = p.outbox.WriteOutboxEvent(ctx, eventbox.Message{
		ID:       uuid.New(),
		Type:     evtypes.PlaceClassCreatedEvent,
		Version:  class.Version,
		Topic:    evtypes.PlacesTopicV1,
		Key:      class.ID.String(),
		Producer: p.identity,
		Payload:  payload,
	})

	return err
}

func (p *Publisher) PublishPlaceClassUpdated(ctx context.Context, class models.PlaceClass) error {
	payload, err := json.Marshal(evtypes.PlaceClassUpdatedPayload{
		PlaceClassID: class.ID,
		ParentID:     class.ParentID,
		Name:         class.Name,
		Description:  class.Description,
		IconKey:      class.IconKey,
		Version:      class.Version,
		UpdatedAt:    class.UpdatedAt,
		DeprecatedAt: class.DeprecatedAt,
	})
	if err != nil {
		return err
	}

	_, err = p.outbox.WriteOutboxEvent(ctx, eventbox.Message{
		ID:       uuid.New(),
		Type:     evtypes.PlaceClassUpdatedEvent,
		Version:  class.Version,
		Topic:    evtypes.PlacesTopicV1,
		Key:      class.ID.String(),
		Producer: p.identity,
		Payload:  payload,
	})

	return err
}
