package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/models"
	"github.com/netbill/places-svc/internal/core/modules/organization"
)

type OrganizationRow struct {
	ID        uuid.UUID `db:"id"`
	Status    string    `db:"status"`
	Name      string    `db:"name"`
	IconKey   *string   `db:"icon_key,omitempty"`
	BannerKey *string   `db:"banner_key,omitempty"`

	Version          int32     `db:"version"`
	SourceCreatedAt  time.Time `db:"source_created_at"`
	SourceUpdatedAt  time.Time `db:"source_updated_at"`
	ReplicaCreatedAt time.Time `db:"replica_created_at"`
	ReplicaUpdatedAt time.Time `db:"replica_updated_at"`
}

func (r OrganizationRow) IsNil() bool {
	return r.ID == uuid.Nil
}

func (r OrganizationRow) ToModel() models.Organization {
	return models.Organization{
		ID:        r.ID,
		Status:    r.Status,
		Name:      r.Name,
		IconKey:   r.IconKey,
		BannerKey: r.BannerKey,
		Version:   r.Version,
		CreatedAt: r.SourceCreatedAt,
		UpdatedAt: r.SourceUpdatedAt,
	}
}

type OrganizationsQ interface {
	New() OrganizationsQ
	Insert(ctx context.Context, input OrganizationRow) error

	Get(ctx context.Context) (OrganizationRow, error)
	Select(ctx context.Context) ([]OrganizationRow, error)

	UpdateOne(ctx context.Context) error

	UpdateStatus(status string) OrganizationsQ
	UpdateName(name string) OrganizationsQ
	UpdateIcon(icon *string) OrganizationsQ
	UpdateBanner(banner *string) OrganizationsQ
	UpdateSourceUpdatedAt(v time.Time) OrganizationsQ
	UpdateVersion(v int32) OrganizationsQ

	FilterByID(id ...uuid.UUID) OrganizationsQ
	FilterByStatus(status string) OrganizationsQ
	FilterByName(name string) OrganizationsQ

	Delete(ctx context.Context) error
}

func (r *Repository) CreateOrganization(
	ctx context.Context,
	input organization.CreateParams,
) error {
	return r.OrganizationsSql.New().Insert(ctx, OrganizationRow{
		ID:              input.ID,
		Status:          input.Status,
		Name:            input.Name,
		IconKey:         input.IconKey,
		BannerKey:       input.BannerKey,
		Version:         1,
		SourceUpdatedAt: input.CreatedAt,
		SourceCreatedAt: input.CreatedAt,
	})
}

func (r *Repository) GetOrganization(
	ctx context.Context,
	orgID uuid.UUID,
) (models.Organization, error) {
	row, err := r.OrganizationsSql.New().FilterByID(orgID).Get(ctx)
	if err != nil {
		return models.Organization{}, nil
	}
	if row.IsNil() {
		return models.Organization{}, errx.ErrorOrganizationNotExists.Raise(
			fmt.Errorf("organization with id %s does not exist", orgID),
		)
	}

	return row.ToModel(), nil
}

func (r *Repository) GetOrgsByIDs(
	ctx context.Context,
	ids []uuid.UUID,
) ([]models.Organization, error) {
	rows, err := r.OrganizationsSql.New().FilterByID(ids...).Select(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]models.Organization, len(rows))
	for i, row := range rows {
		res[i] = row.ToModel()
	}

	return res, nil
}

func (r *Repository) UpdateOrganization(
	ctx context.Context,
	orgID uuid.UUID,
	params organization.UpdateParams,
) error {
	return r.OrganizationsSql.New().
		FilterByID(orgID).
		UpdateName(params.Name).
		UpdateIcon(params.IconKey).
		UpdateBanner(params.BannerKey).
		UpdateSourceUpdatedAt(params.UpdatedAt).
		UpdateVersion(params.Version).
		UpdateOne(ctx)
}

func (r *Repository) DeleteOrganization(ctx context.Context, ID uuid.UUID) error {
	return r.OrganizationsSql.New().
		FilterByID(ID).
		Delete(ctx)
}
