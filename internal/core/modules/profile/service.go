package profile

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/models"
)

type Service struct {
	repo repo
}

type repo interface {
	UpsertProfile(ctx context.Context, profile models.Profile) error
	UpdateProfileUsername(ctx context.Context, accountID uuid.UUID, username string) error
	DeleteProfile(ctx context.Context, ID uuid.UUID) error

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

func New(repo repo) Service {
	return Service{repo: repo}
}
