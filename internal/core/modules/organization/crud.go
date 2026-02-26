package organization

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/models"
)

type CreateParams struct {
	ID        uuid.UUID `json:"id"`
	Status    string    `json:"status"`
	Name      string    `json:"name"`
	IconKey   *string   `json:"icon_key,omitempty"`
	BannerKey *string   `json:"banner_key,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}

func (m *Module) Create(
	ctx context.Context,
	org CreateParams,
) error {
	return m.repo.CreateOrganization(ctx, org)
}

func (m *Module) Get(
	ctx context.Context,
	orgID uuid.UUID,
) (models.Organization, error) {
	return m.repo.GetOrganization(ctx, orgID)
}

func (m *Module) GetByIDs(
	ctx context.Context,
	ids []uuid.UUID,
) ([]models.Organization, error) {
	res, err := m.repo.GetOrgsByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	return res, nil
}

type UpdateParams struct {
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	IconKey   *string   `json:"icon_key"`
	BannerKey *string   `json:"banner_key"`
	Version   int32     `json:"version"`
	UpdatedAt time.Time `json:"updated_at"`
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

	return m.repo.UpdateOrganization(ctx, orgID, params)
}

func (m *Module) Delete(
	ctx context.Context,
	organizationID uuid.UUID,
) error {
	return m.repo.DeleteOrganization(ctx, organizationID)
}
