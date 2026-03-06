package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/organization"
	"github.com/netbill/places-svc/internal/errx"
	"github.com/netbill/places-svc/internal/models"
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
	Exists(ctx context.Context) (bool, error)

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

type OrgRepository struct {
	query OrganizationsQ
}

func NewOrgRepository(query OrganizationsQ) *OrgRepository {
	return &OrgRepository{
		query: query,
	}
}

func (r *OrgRepository) Create(
	ctx context.Context,
	input organization.CreateParams,
) error {
	return r.query.New().Insert(ctx, OrganizationRow{
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

func (r *OrgRepository) Get(
	ctx context.Context,
	orgID uuid.UUID,
) (models.Organization, error) {
	row, err := r.query.New().FilterByID(orgID).Get(ctx)
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

func (r *OrgRepository) GetListByIDs(
	ctx context.Context,
	ids []uuid.UUID,
) ([]models.Organization, error) {
	rows, err := r.query.New().FilterByID(ids...).Select(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]models.Organization, len(rows))
	for i, row := range rows {
		res[i] = row.ToModel()
	}

	return res, nil
}

func (r *OrgRepository) Exists(ctx context.Context, orgID uuid.UUID) (bool, error) {
	return r.query.New().FilterByID(orgID).Exists(ctx)
}

func (r *OrgRepository) Update(
	ctx context.Context,
	orgID uuid.UUID,
	params organization.UpdateParams,
) error {
	return r.query.New().
		FilterByID(orgID).
		UpdateName(params.Name).
		UpdateIcon(params.IconKey).
		UpdateBanner(params.BannerKey).
		UpdateSourceUpdatedAt(params.UpdatedAt).
		UpdateVersion(params.Version).
		UpdateOne(ctx)
}

func (r *OrgRepository) Delete(ctx context.Context, ID uuid.UUID) error {
	return r.query.New().
		FilterByID(ID).
		Delete(ctx)
}
