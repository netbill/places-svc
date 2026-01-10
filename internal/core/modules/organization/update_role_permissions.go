package organization

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
)

func (s Service) UpdateRolePermissions(
	ctx context.Context,
	roleID uuid.UUID,
	permissions map[string]bool,
) error {
	if err := s.repo.UpdateRolePermissions(ctx, roleID, permissions); err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to set role permissions: %w", err),
		)
	}

	return nil
}
