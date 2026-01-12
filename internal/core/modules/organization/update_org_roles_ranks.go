package organization

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
)

func (s Service) UpdateOrgRolesRanks(
	ctx context.Context,
	organizationID uuid.UUID,
	order map[uuid.UUID]uint,
) error {
	if err := s.repo.UpdateOrgRolesRanks(ctx, organizationID, order); err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update roles ranks: %w", err),
		)
	}

	return nil
}
