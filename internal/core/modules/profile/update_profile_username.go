package profile

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
)

func (s Service) UpdateProfileUsername(ctx context.Context, accountID uuid.UUID, username string) error {
	if err := s.repo.UpdateProfileUsername(ctx, accountID, username); err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to upsert profile: %w", err),
		)
	}

	return nil
}
