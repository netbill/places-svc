package profile

import (
	"context"
	"fmt"

	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/models"
)

func (s Service) UpsertProfile(ctx context.Context, profile models.Profile) error {
	if err := s.repo.UpsertProfile(ctx, profile); err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to upsert profile: %w", err),
		)
	}

	return nil
}
