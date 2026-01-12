package organization

import (
	"context"
	"fmt"

	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/models"
)

func (s Service) UpsertOrgRole(
	ctx context.Context,
	params models.OrgRole,
) error {
	if err := s.repo.UpsertOrgRole(ctx, params); err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to create role: %w", err),
		)
	}

	return nil
}
