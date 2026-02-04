package profile

import (
	"context"

	"github.com/google/uuid"
)

func (m *Module) DeleteProfile(
	ctx context.Context,
	accountID uuid.UUID,
) error {
	return m.repo.DeleteProfile(ctx, accountID)
}
