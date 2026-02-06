package organization

import (
	"context"

	"github.com/netbill/places-svc/internal/core/models"
)

func (m *Module) CreateOrgMember(ctx context.Context, member models.OrgMember) (models.OrgMember, error) {
	return m.repo.CreateOrgMember(ctx, member)
}
