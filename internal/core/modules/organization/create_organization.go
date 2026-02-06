package organization

import (
	"context"

	"github.com/netbill/places-svc/internal/core/models"
)

func (m *Module) CreateOrganization(
	ctx context.Context,
	org models.Organization,
) (models.Organization, error) {
	return m.repo.CreateOrganization(ctx, org)
}
