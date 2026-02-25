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
	Label          *string   `json:"label,omitempty"`
	Position       *string   `json:"position,omitempty"`

	Version   *string   `json:"version,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
