package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

const (
	PlaceStatusActive    = "active"
	PlaceStatusInactive  = "inactive"
	PlaceStatusSuspended = "suspended"
)

type Place struct {
	ID             uuid.UUID  `json:"id"`
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

func (p Place) IsNil() bool {
	return p.ID == uuid.Nil
}
