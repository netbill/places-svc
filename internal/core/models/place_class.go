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

type PlaceClassMedia struct {
	Icon *string `json:"icon,omitempty"`
}

type PlaceClassUploadMediaLinks struct {
	IconUploadURL string `json:"icon_upload_url"`
	IconGetURL    string `json:"icon_get_url"`
}

type UpdatePlaceClassMedia struct {
	Links           PlaceClassUploadMediaLinks `json:"links"`
	UploadSessionID uuid.UUID                  `json:"upload_sessian_id"`
	UploadToken     string                     `json:"upload_token"`
}
