package models

import (
	"time"

	"github.com/google/uuid"
)

type PlaceClass struct {
	ID       uuid.UUID  `json:"id"`
	ParentID *uuid.UUID `json:"parent_id,omitempty"`

	Name        string  `json:"name"`
	Description string  `json:"description"`
	IconKey     *string `json:"icon_key,omitempty"`

	Version      int32      `json:"version"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeprecatedAt *time.Time `json:"deprecated_at,omitempty"`
}

type UploadPlaceClassMediaLinks struct {
	Icon UploadMediaLink `json:"icon"`
}
