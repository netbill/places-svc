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

type OrgMemberRow struct {
	ID               uuid.UUID `db:"id"`
	AccountID        uuid.UUID `db:"account_id"`
	OrganizationID   uuid.UUID `db:"organization_id"`
	Head             bool      `db:"head"`
	Label            *string   `db:"label,omitempty"`
	Position         *string   `db:"position,omitempty"`
	Version          int32     `db:"version"`
	SourceCreatedAt  time.Time `db:"source_created_at"`
	SourceUpdatedAt  time.Time `db:"source_updated_at"`
	ReplicaCreatedAt time.Time `db:"replica_created_at"`
	ReplicaUpdatedAt time.Time `db:"replica_updated_at"`
}

func (r OrgMemberRow) IsNil() bool {
	return r.ID == uuid.Nil
}

func (r OrgMemberRow) ToModel() models.OrgMember {
	return models.OrgMember{
		ID:             r.ID,
		AccountID:      r.AccountID,
		OrganizationID: r.OrganizationID,
		Head:           r.Head,
		Label:          r.Label,
		Position:       r.Position,
		Version:        r.Version,
		CreatedAt:      r.SourceCreatedAt,
		UpdatedAt:      r.SourceUpdatedAt,
	}
}

type OrgMembersQ interface {
	New() OrgMembersQ
	Insert(ctx context.Context, input OrgMemberRow) error

	Get(ctx context.Context) (OrgMemberRow, error)
	Select(ctx context.Context) ([]OrgMemberRow, error)
	Exists(ctx context.Context) (bool, error)

	FilterByID(id ...uuid.UUID) OrgMembersQ
	FilterByAccountID(accountID ...uuid.UUID) OrgMembersQ
	FilterByOrganizationID(organizationID ...uuid.UUID) OrgMembersQ
	FilterByHead(head bool) OrgMembersQ

	UpdateOne(ctx context.Context) error

	UpdateVersion(v int32) OrgMembersQ
	UpdateLabel(label *string) OrgMembersQ
	UpdatePosition(position *string) OrgMembersQ
	UpdateSourceUpdatedAt(sourceUpdatedAt time.Time) OrgMembersQ

	Delete(ctx context.Context) error
}

type OrgMembersRepository struct {
	query OrgMembersQ
}

func NewOrgMember(query OrgMembersQ) *OrgMembersRepository {
	return &OrgMembersRepository{
		query: query,
	}
}

func (r *OrgMembersRepository) Create(
	ctx context.Context,
	params organization.CreateMemberParams,
) error {
	return r.query.New().Insert(ctx, OrgMemberRow{
		ID:              params.ID,
		AccountID:       params.AccountID,
		OrganizationID:  params.OrganizationID,
		Label:           params.Label,
		Position:        params.Position,
		Head:            params.Head,
		Version:         1,
		SourceCreatedAt: params.CreatedAt,
		SourceUpdatedAt: params.CreatedAt,
	})
}

func (r *OrgMembersRepository) Update(
	ctx context.Context,
	memberID uuid.UUID,
	params organization.UpdateMemberParams,
) error {
	return r.query.New().FilterByID(memberID).
		UpdateLabel(params.Label).
		UpdatePosition(params.Position).
		UpdateSourceUpdatedAt(params.UpdatedAt).
		UpdateVersion(params.Version).
		UpdateOne(ctx)
}

func (r *OrgMembersRepository) Delete(ctx context.Context, memberID uuid.UUID) error {
	return r.query.New().
		FilterByID(memberID).
		Delete(ctx)
}

func (r *OrgMembersRepository) GetForAccountAndOrg(
	ctx context.Context,
	accountID, organizationID uuid.UUID,
) (models.OrgMember, error) {
	row, err := r.query.New().
		FilterByOrganizationID(organizationID).
		FilterByAccountID(accountID).
		Get(ctx)
	if err != nil {
		return models.OrgMember{}, err
	}

	return row.ToModel(), nil
}

func (r *OrgMembersRepository) ExistsForAccountAndOrg(
	ctx context.Context,
	accountID, organizationID uuid.UUID,
) (bool, error) {
	return r.query.New().
		FilterByOrganizationID(organizationID).
		FilterByAccountID(accountID).
		Exists(ctx)
}

func (r *OrgMembersRepository) GetByID(ctx context.Context, memberID uuid.UUID) (models.OrgMember, error) {
	row, err := r.query.New().FilterByID(memberID).Get(ctx)
	if err != nil {
		return models.OrgMember{}, err
	}
	if row.IsNil() {
		return models.OrgMember{}, errx.ErrorOrgMemberNotExists.Raise(
			fmt.Errorf("organization member with id %s not found", memberID),
		)
	}

	return row.ToModel(), nil
}

func (r *OrgMembersRepository) ExistsByID(ctx context.Context, memberID uuid.UUID) (bool, error) {
	return r.query.New().FilterByID(memberID).Exists(ctx)
}
