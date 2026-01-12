package organization

import (
	"context"
	"fmt"

	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/models"
)

func (s Service) UpdateOrganizationStatus(ctx context.Context, org models.Organization) error {
	err := s.repo.UpsertOrganization(ctx, org)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to activate organization: %w", err))
	}

	return nil
}
