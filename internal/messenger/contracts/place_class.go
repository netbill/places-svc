package contracts

import (
	"time"

	"github.com/google/uuid"
)

const PlaceClassesTopicV1 = "place.classes.v1"

const PlaceClassCreatedEvent = "place_class.created"

type PlaceClassCreatedPayload struct {
	PlaceClassID uuid.UUID  `json:"place_class_id"`
	ParentID     *uuid.UUID `json:"parent_id,omitempty"`
	Code         string     `json:"code"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	Icon         *string    `json:"icon"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

const PlaceClassUpdatedEvent = "place_class.updated"

type PlaceClassUpdatedPayload struct {
	PlaceClassID uuid.UUID  `json:"place_class_id"`
	ParentID     *uuid.UUID `json:"parent_id,omitempty"`
	Code         string     `json:"code"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	Icon         *string    `json:"icon"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

const PlaceClassParentUpdatedEvent = "place_class.parent.updated"

type PlaceClassParentUpdatedPayload struct {
	PlaceClassID uuid.UUID  `json:"place_class_id"`
	ParentID     *uuid.UUID `json:"parent_id,omitempty"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

const PlaceClassDeletedEvent = "place_class.deleted"

type PlaceClassDeletedPayload struct {
	PlaceClassID uuid.UUID `json:"place_class_id"`
	DeletedAt    time.Time `json:"deleted_at"`
}

const PlacesClassReplacedEvent = "places_class.replaced"

type PlacesClassReplacedPayload struct {
	OldClassID uuid.UUID `json:"old_class_id"`
	NewClassID uuid.UUID `json:"new_class_id"`
	ReplacedAt time.Time `json:"replaced_at"`
}
