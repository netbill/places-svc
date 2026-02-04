package organization

import (
	"context"

	"github.com/google/uuid"
)

func (m *Module) DeleteOrganization(
	ctx context.Context,
	organizationID uuid.UUID,
) error {
	return m.repo.DeleteOrganization(ctx, organizationID)
}
