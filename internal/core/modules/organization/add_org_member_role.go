package organization

import (
	"context"
	"time"

	"github.com/google/uuid"
)

func (m *Module) AddOrgMemberRole(
	ctx context.Context,
	memberID, roleID uuid.UUID,
	addedAt time.Time,
) error {
	return m.repo.AddOrgMemberRole(ctx, memberID, roleID, addedAt)
}
