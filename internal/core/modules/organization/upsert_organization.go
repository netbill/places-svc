package organization

import (
	"context"
	"fmt"

	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/models"
)

func (s Service) UpsertOrganization(
	ctx context.Context,
	org models.Organization,
) error {
	if err := s.repo.UpsertOrganization(ctx, org); err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to create organization: %w", err),
		)
	}

	return nil
}
