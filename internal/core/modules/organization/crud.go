package organization

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/models"
)

func (m *Module) Create(
	ctx context.Context,
	org models.Organization,
) error {
	return m.repo.Create(ctx, org)
}

func (m *Module) Get(
	ctx context.Context,
	orgID uuid.UUID,
) (models.Organization, error) {
	return m.repo.GetOrganization(ctx, orgID)
}

type UpdateParams struct {
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	IconKey   *string   `json:"icon_key"`
	BannerKey *string   `json:"banner_key"`
	Version   int32     `json:"version"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (m *Module) Update(
	ctx context.Context,
	orgID uuid.UUID,
	params UpdateParams,
) error {
	org, err := m.repo.GetOrganization(ctx, orgID)
	if err != nil {
		return err
	}

	if org.Version >= params.Version {
		return nil // No update needed
	}

	return m.repo.Update(ctx, orgID, params)
}

func (m *Module) Delete(
	ctx context.Context,
	organizationID uuid.UUID,
) error {
	return m.repo.Delete(ctx, organizationID)
}
