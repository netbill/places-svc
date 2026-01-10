package models

import (
	"time"

	"github.com/google/uuid"
)

type Member struct {
	ID             uuid.UUID `json:"id"`
	AccountID      uuid.UUID `json:"account_id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (m Member) IsNil() bool {
	return m.ID == uuid.Nil
}
