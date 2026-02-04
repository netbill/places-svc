package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/models"
)

type OrgRoleRow struct {
	ID             uuid.UUID `db:"id"`
	OrganizationID uuid.UUID `db:"organization_id"`
	Rank           uint      `db:"rank"`

	SourceCreatedAt  time.Time `db:"source_created_at"`
	SourceUpdatedAt  time.Time `db:"source_updated_at"`
	ReplicaCreatedAt time.Time `db:"replica_created_at"`
	ReplicaUpdatedAt time.Time `db:"replica_updated_at"`
}

func (r OrgRoleRow) IsNil() bool {
	return r.ID == uuid.Nil
}

func (r OrgRoleRow) ToModel() models.OrgRole {
	return models.OrgRole{
		ID:             r.ID,
		OrganizationID: r.OrganizationID,
		Rank:           r.Rank,
		CreatedAt:      r.SourceCreatedAt,
		UpdatedAt:      r.SourceUpdatedAt,
	}
}

type OrgRolesQ interface {
	New() OrgRolesQ
	Insert(ctx context.Context, input OrgRoleRow) (OrgRoleRow, error)

	Get(ctx context.Context) (OrgRoleRow, error)
	Select(ctx context.Context) ([]OrgRoleRow, error)

	UpdateMany(ctx context.Context) (int64, error)
	UpdateOne(ctx context.Context) (OrgRoleRow, error)

	UpdateRank(rank uint) OrgRolesQ
	UpdateSourceUpdatedAt(v time.Time) OrgRolesQ

	FilterByID(id ...uuid.UUID) OrgRolesQ
	FilterByOrganizationID(organizationID ...uuid.UUID) OrgRolesQ
	FilterByRank(rank uint) OrgRolesQ

	Delete(ctx context.Context) error
}
