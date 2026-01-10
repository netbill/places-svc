package organization

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
)

func (s Service) DeleteRole(ctx context.Context, roleID uuid.UUID) error {
	if err := s.repo.DeleteRole(ctx, roleID); err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to delete role: %w", err),
		)
	}

	return nil
}
