package contracts

import (
	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/models"
)

const PlaceClassCreatedEvent = "place_class.created"

type PlaceClassCreatedPayload struct {
	PlaceClass models.PlaceClass `json:"place_class"`
}

const PlaceClassUpdatedEvent = "place_class.updated"

type PlaceClassUpdatedPayload struct {
	PlaceClass models.PlaceClass `json:"place_class"`
}

const PlaceClassParentUpdatedEvent = "place_class.parent.updated"

type PlaceClassParentUpdatedPayload struct {
	PlaceClass models.PlaceClass `json:"place_class"`
}

const PlaceClassDeletedEvent = "place_class.deleted"

type PlaceClassDeletedPayload struct {
	PlaceClassID uuid.UUID `json:"place_class_id"`
}

const PlacesClassReplacedEvent = "places_class.replaced"

type PlacesClassReplacedPayload struct {
	OldClassID uuid.UUID `json:"old_class_id"`
	NewClassID uuid.UUID `json:"new_class_id"`
}
