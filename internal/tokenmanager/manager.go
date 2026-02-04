package tokenmanager

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/netbill/restkit/tokens"
)

type Manager struct {
	uploadSK string

	placeMediaTokenTTL      time.Duration
	placeClassMediaTokenTTL time.Duration
}

const (
	PlacesActor        = "places-svc"
	PlaceResource      = "place"
	PlaceClassResource = "place-class"
)

func New(uploadSK string, placeMediaTokenTTL, placeClassMediaTokenTTL time.Duration) *Manager {
	return &Manager{
		uploadSK:                uploadSK,
		placeMediaTokenTTL:      placeMediaTokenTTL,
		placeClassMediaTokenTTL: placeClassMediaTokenTTL,
	}
}

func (m *Manager) NewUploadPlaceMediaToken(
	OwnerAccountID uuid.UUID,
	ResourceID uuid.UUID,
	UploadSessionID uuid.UUID,
) (string, error) {
	tkn, err := tokens.UploadContentClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   OwnerAccountID.String(),
			Issuer:    PlacesActor,
			Audience:  []string{PlacesActor},
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(m.placeMediaTokenTTL)),
		},
		UploadSessionID: UploadSessionID,
		ResourceID:      ResourceID.String(),
		Resource:        PlaceResource,
	}.GenerateJWT(m.uploadSK)
	if err != nil {
		return "", fmt.Errorf("failed to generate upload place media token, cause: %w", err)
	}

	return tkn, nil
}

func (m *Manager) NewUploadPlaceClassMediaToken(
	OwnerAccountID uuid.UUID,
	ResourceID uuid.UUID,
	UploadSessionID uuid.UUID,
) (string, error) {
	tkn, err := tokens.UploadContentClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   OwnerAccountID.String(),
			Issuer:    PlacesActor,
			Audience:  []string{PlacesActor},
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(m.placeClassMediaTokenTTL)),
		},
		UploadSessionID: UploadSessionID,
		ResourceID:      ResourceID.String(),
		Resource:        PlaceClassResource,
	}.GenerateJWT(m.uploadSK)
	if err != nil {
		return "", fmt.Errorf("failed to generate upload place class media token, cause: %w", err)
	}

	return tkn, nil
}
