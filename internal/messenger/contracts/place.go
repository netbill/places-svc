package contracts

import (
	"time"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

const PlacesTopicV1 = "places.v1"

const PlaceCreatedEvent = "place.created"

type PlaceCreatedPayload struct {
	PlaceID        uuid.UUID  `json:"place_id"`
	ClassID        uuid.UUID  `json:"class_id"`
	OrganizationID *uuid.UUID `json:"organization_id,omitempty"`

	Status   string    `json:"status"`
	Verified bool      `json:"verified"`
	Point    orb.Point `json:"point,omitempty"`
	Address  string    `json:"address,omitempty"`
	Name     string    `json:"name"`

	Description *string `json:"description"`
	Icon        *string `json:"icon"`
	Banner      *string `json:"banner"`
	Website     *string `json:"website"`
	Phone       *string `json:"phone"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

const PlaceUpdatedEvent = "place.updated"

type PlaceUpdatedPayload struct {
	PlaceID        uuid.UUID  `json:"place_id"`
	ClassID        uuid.UUID  `json:"class_id"`
	OrganizationID *uuid.UUID `json:"organization_id,omitempty"`

	Status   string `json:"status"`
	Verified bool   `json:"verified"`
	Address  string `json:"address,omitempty"`
	Name     string `json:"name"`

	Description *string   `json:"description"`
	Icon        *string   `json:"icon"`
	Banner      *string   `json:"banner"`
	Website     *string   `json:"website"`
	Phone       *string   `json:"phone"`
	UpdatedAt   time.Time `json:"updated_at"`
}

const PlaceStatusUpdatedEvent = "place.status.updated"

type PlaceStatusUpdatedPayload struct {
	PlaceID   uuid.UUID `json:"place_id"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
}

const PlaceVerifiedUpdatedEvent = "place.verified.updated"

type PlaceVerifiedUpdatedPayload struct {
	PlaceID   uuid.UUID `json:"place_id"`
	Verified  bool      `json:"verified"`
	UpdatedAt time.Time `json:"updated_at"`
}

const PlaceClassIDUpdatedEvent = "place.class_id.updated"

type PlaceClassIDUpdatedPayload struct {
	PlaceID   uuid.UUID `json:"place_id"`
	ClassID   uuid.UUID `json:"class_id"`
	UpdatedAt time.Time `json:"updated_at"`
}

const PlaceDeletedEvent = "place.deleted"

type PlaceDeletedPayload struct {
	PlaceID   uuid.UUID `json:"place_id"`
	DeletedAt time.Time `json:"deleted_at"`
}
