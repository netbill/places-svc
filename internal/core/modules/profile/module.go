package profile

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/models"
)

type Module struct {
	repo repo
}

type repo interface {
	CreateProfile(ctx context.Context, profile models.Profile) error
	UpdateProfile(ctx context.Context, profileID uuid.UUID, params UpdateParams) error
	DeleteProfile(ctx context.Context, ID uuid.UUID) error

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

func New(repo repo) *Module {
	return &Module{repo: repo}
}
