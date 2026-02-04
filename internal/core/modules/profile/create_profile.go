package profile

import (
	"context"

	"github.com/netbill/places-svc/internal/core/models"
)

func (m *Module) CreateProfile(ctx context.Context, profile models.Profile) error {
	return m.repo.CreateProfile(ctx, profile)
}
