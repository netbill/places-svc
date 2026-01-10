package organization

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
)

func (s Service) DeleteOrganization(ctx context.Context, ID uuid.UUID) error {
	if err := s.repo.DeleteOrganization(ctx, ID); err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to delete organization: %w", err),
		)
	}

	return nil
}
