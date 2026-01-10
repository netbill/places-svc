package models

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Head           bool      `json:"head"`
	Rank           uint      `json:"rank"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (r Role) IsNil() bool {
	return r.ID == uuid.Nil
}
