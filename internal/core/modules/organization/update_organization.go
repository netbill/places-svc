package organization

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/models"
)

type UpdateParams struct {
	Name      string    `json:"name"`
	Icon      *string   `json:"icon"`
	Banner    *string   `json:"banner"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (m *Module) UpdateOrganization(
	ctx context.Context,
	orgID uuid.UUID,
	params UpdateParams,
) (models.Organization, error) {
	return m.repo.UpdateOrganization(ctx, orgID, params)
}
