package models

import (
	"time"

	"github.com/google/uuid"
)

type OrgMember struct {
	ID             uuid.UUID `json:"id"`
	AccountID      uuid.UUID `json:"account_id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Head           bool      `json:"head"`
	Position       *string   `json:"position,omitempty"`
	Label          *string   `json:"label,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
