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

type PlaceMedia struct {
	Icon   *string `json:"icon,omitempty"`
	Banner *string `json:"banner,omitempty"`
}

type PlaceUploadMediaLinks struct {
	IconUploadURL   string `json:"icon_upload_url"`
	IconGetURL      string `json:"icon_get_url"`
	BannerUploadURL string `json:"banner_upload_url"`
	BannerGetURL    string `json:"banner_get_url"`
}

type UpdatePlaceMedia struct {
	Links           PlaceUploadMediaLinks `json:"links"`
	UploadSessionID uuid.UUID             `json:"upload_session_id"`
	UploadToken     string                `json:"upload_token"`
}
