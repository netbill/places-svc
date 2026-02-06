package organization

import (
	"context"

	"github.com/netbill/places-svc/internal/core/models"
)

func (m *Module) CreateOrgRole(ctx context.Context, role models.OrgRole) (models.OrgRole, error) {
	return m.repo.CreateOrgRole(ctx, role)
}
