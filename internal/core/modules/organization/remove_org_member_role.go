package organization

import (
	"context"

	"github.com/google/uuid"
)

func (m *Module) RemoveOrgMemberRole(
	ctx context.Context,
	memberID, roleID uuid.UUID,
) error {
	return m.repo.RemoveOrgMemberRole(ctx, memberID, roleID)
}
