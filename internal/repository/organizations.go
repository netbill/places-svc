package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/models"
)

type OrganizationRow struct {
	ID       uuid.UUID `db:"id"`
	Status   string    `db:"status"`
	Verified bool      `db:"verified"`
	Name     string    `db:"name"`
	Icon     *string   `db:"icon"`
	Banner   *string   `db:"banner,omitempty"`

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
		Verified:  r.Verified,
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
	UpdateVerified(verified bool) OrganizationsQ
	UpdateName(name string) OrganizationsQ
	UpdateIcon(icon *string) OrganizationsQ
	UpdateBanner(banner *string) OrganizationsQ
	UpdateSourceUpdatedAt(v time.Time) OrganizationsQ

	FilterByID(id ...uuid.UUID) OrganizationsQ
	FilterByStatus(status string) OrganizationsQ
	FilterByVerified(verified bool) OrganizationsQ
	FilterByName(name string) OrganizationsQ

	Delete(ctx context.Context) error
}
