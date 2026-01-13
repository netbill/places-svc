package models

import (
	"time"

	"github.com/google/uuid"
)

type PlaceClass struct {
	ID          uuid.UUID  `json:"id"`
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
	Code        string     `json:"code"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Icon        *string    `json:"icon"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (c PlaceClass) IsNil() bool {
	return c.ID == uuid.Nil
}
