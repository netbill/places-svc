package organization

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
)

func (s Service) DeleteMember(ctx context.Context, ID uuid.UUID) error {
	if err := s.repo.DeleteMember(ctx, ID); err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to delete member %s: %w", ID, err),
		)
	}

	return nil
}
