package organization

import (
	"context"

	"github.com/google/uuid"
)

func (m *Module) DeleteOrgRole(
	ctx context.Context,
	roleID uuid.UUID,
) error {
	return m.repo.DeleteOrgRole(ctx, roleID)
}
