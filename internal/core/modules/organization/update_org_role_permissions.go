package organization

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
)

func (s Service) UpdateOrgRolePermissions(
	ctx context.Context,
	roleID uuid.UUID,
	permissions map[string]bool,
) error {
	if err := s.repo.UpdateOrgRolePermissions(ctx, roleID, permissions); err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to set role permissions: %w", err),
		)
	}

	return nil
}
