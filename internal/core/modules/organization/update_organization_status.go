package organization

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/models"
)

func (s Service) UpdateOrganizationStatus(ctx context.Context, ID uuid.UUID, status string) (models.Organization, error) {
	org, err := s.repo.UpdateOrganizationStatus(ctx, ID, status)
	if err != nil {
		return models.Organization{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to activate organization: %w", err))
	}

	return org, nil
}
