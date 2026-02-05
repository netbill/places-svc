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
	ID     uuid.UUID `db:"id"`
	Status string    `db:"status"`
	Name   string    `db:"name"`
	Icon   *string   `db:"icon"`
	Banner *string   `db:"banner,omitempty"`

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
		Icon:      r.Icon,
		Banner:    r.Banner,
		CreatedAt: r.SourceCreatedAt,
		UpdatedAt: r.SourceUpdatedAt,
	}
}

type OrganizationsQ interface {
	New() OrganizationsQ
	Insert(ctx context.Context, input OrganizationRow) (OrganizationRow, error)

	Get(ctx context.Context) (OrganizationRow, error)
	Select(ctx context.Context) ([]OrganizationRow, error)

	UpdateMany(ctx context.Context) (int64, error)
	UpdateOne(ctx context.Context) (OrganizationRow, error)

	UpdateStatus(status string) OrganizationsQ
	UpdateName(name string) OrganizationsQ
	UpdateIcon(icon *string) OrganizationsQ
	UpdateBanner(banner *string) OrganizationsQ
	UpdateSourceUpdatedAt(v time.Time) OrganizationsQ

	FilterByID(id ...uuid.UUID) OrganizationsQ
	FilterByStatus(status string) OrganizationsQ
	FilterByName(name string) OrganizationsQ

	Delete(ctx context.Context) error
}

func (r *Repository) CreateOrganization(
	ctx context.Context,
	input models.Organization,
) (models.Organization, error) {
	res, err := r.OrganizationsQ.New().Insert(ctx, OrganizationRow{
		ID:              input.ID,
		Status:          input.Status,
		Name:            input.Name,
		Icon:            input.Icon,
		Banner:          input.Banner,
		SourceCreatedAt: input.CreatedAt,
		SourceUpdatedAt: input.UpdatedAt,
	})
	if err != nil {
		return models.Organization{}, fmt.Errorf(
			"failed to creating organization, cause: %w", err,
		)
	}

	return res.ToModel(), err
}

func (r *Repository) UpdateOrganization(
	ctx context.Context,
	orgID uuid.UUID,
	params organization.UpdateParams,
) (models.Organization, error) {
	res, err := r.OrganizationsQ.New().
		FilterByID(orgID).
		UpdateName(params.Name).
		UpdateIcon(params.Icon).
		UpdateBanner(params.Banner).
		UpdateSourceUpdatedAt(params.UpdatedAt).
		UpdateOne(ctx)
	if err != nil {
		return models.Organization{}, fmt.Errorf(
			"failed to update organization, cause: %s", err,
		)
	}

	if res.IsNil() {
		return models.Organization{}, errx.ErrorOrganizationNotExists.Raise(
			fmt.Errorf("organization with ID %s not found", orgID),
		)
	}

	return res.ToModel(), err
}

func (r *Repository) UpdateOrgStatus(
	ctx context.Context,
	orgID uuid.UUID,
	status string,
	updatedAt time.Time,
) (models.Organization, error) {
	res, err := r.OrganizationsQ.New().
		FilterByID(orgID).
		UpdateStatus(status).
		UpdateSourceUpdatedAt(updatedAt).
		UpdateOne(ctx)
	if err != nil {
		return models.Organization{}, err
	}

	if res.IsNil() {
		return models.Organization{}, errx.ErrorOrganizationNotExists.Raise(
			fmt.Errorf("organization with ID %s not found", orgID),
		)
	}

	return res.ToModel(), err
}

func (r *Repository) DeleteOrganization(ctx context.Context, ID uuid.UUID) error {
	return r.OrganizationsQ.New().
		FilterByID(ID).
		Delete(ctx)
}
