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

	OrderByRank(asc bool) OrgRolesQ

	UpdateRolesRanks(
		ctx context.Context,
		organizationID uuid.UUID,
		order map[uuid.UUID]uint,
		updatedAt time.Time,
	) error

	Delete(ctx context.Context) error
}

func (r *Repository) CreateOrgRole(
	ctx context.Context,
	role models.OrgRole,
) (models.OrgRole, error) {
	row, err := r.OrgRolesQ.New().Insert(ctx, OrgRoleRow{
		ID:             role.ID,
		OrganizationID: role.OrganizationID,
		Rank:           role.Rank,
	})
	if err != nil {
		return models.OrgRole{}, err
	}
	return row.ToModel(), nil
}

func (r *Repository) UpdateOrgRole(
	ctx context.Context,
	role models.OrgRole,
) (models.OrgRole, error) {
	row, err := r.OrgRolesQ.New().
		FilterByID(role.ID).
		UpdateOne(ctx)
	if err != nil {
		return models.OrgRole{}, err
	}
	return row.ToModel(), nil
}

func (r *Repository) UpdateOrgRolesRanks(
	ctx context.Context,
	organizationID uuid.UUID,
	order map[uuid.UUID]uint,
	updatedAt time.Time,
) error {
	return r.OrgRolesQ.New().UpdateRolesRanks(
		ctx,
		organizationID,
		order,
		updatedAt,
	)
}

func (r *Repository) UpdateOrgRolePermissions(
	ctx context.Context,
	roleID uuid.UUID,
	permissions map[string]bool,
	updatedAt time.Time,
) error {
	//TODO: implement

	return nil
}

func (r *Repository) DeleteOrgRole(
	ctx context.Context,
	roleID uuid.UUID,
) error {
	return r.OrgRolesQ.New().
		FilterByID(roleID).
		Delete(ctx)
}
