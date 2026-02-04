package organization

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type UpdateParams struct {
	Verified  bool      `json:"verified"`
	Name      string    `json:"name"`
	Icon      *string   `json:"icon"`
	Banner    *string   `json:"banner"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (m *Module) UpdateOrganization(
	ctx context.Context,
	orgID uuid.UUID,
	params UpdateParams,
) error {
	return m.repo.UpdateOrganization(ctx, orgID, params)
}
