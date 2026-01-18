package contracts

import (
	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/models"
)

const PlaceCreatedEvent = "place.created"

type PlaceCreatedPayload struct {
	Place models.Place `json:"place"`
}

const PlaceUpdatedEvent = "place.updated"

type PlaceUpdatedPayload struct {
	Place models.Place `json:"place"`
}

const PlaceStatusUpdatedEvent = "place.status.updated"

type PlaceStatusUpdatedPayload struct {
	Place models.Place `json:"place"`
}

const PlaceVerifiedUpdatedEvent = "place.verified.updated"

type PlaceVerifiedUpdatedPayload struct {
	Place models.Place `json:"place"`
}

const PlaceClassIDUpdatedEvent = "place.class_id.updated"

type PlaceClassIDUpdatedPayload struct {
	Place models.Place `json:"place"`
}

const PlaceDeletedEvent = "place.deleted"

type PlaceDeletedPayload struct {
	PlaceID uuid.UUID `json:"place_id"`
}
