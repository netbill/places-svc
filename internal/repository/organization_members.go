package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/models"
)

type OrgMemberRow struct {
	ID               uuid.UUID `db:"id"`
	AccountID        uuid.UUID `db:"account_id"`
	OrganizationID   uuid.UUID `db:"organization_id"`
	Head             bool      `db:"head"`
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
		CreatedAt:      r.SourceCreatedAt,
		UpdatedAt:      r.SourceUpdatedAt,
	}
}

type OrgMembersQ interface {
	New() OrgMembersQ
	Insert(ctx context.Context, input OrgMemberRow) (OrgMemberRow, error)

	Get(ctx context.Context) (OrgMemberRow, error)
	Select(ctx context.Context) ([]OrgMemberRow, error)

	FilterByID(id ...uuid.UUID) OrgMembersQ
	FilterByAccountID(accountID ...uuid.UUID) OrgMembersQ
	FilterByOrganizationID(organizationID ...uuid.UUID) OrgMembersQ
	FilterByHead(head bool) OrgMembersQ

	Delete(ctx context.Context) error
}

func (r *Repository) CreateOrgMember(
	ctx context.Context,
	member models.OrgMember,
) (models.OrgMember, error) {
	row, err := r.OrgMembersQ.New().Insert(ctx, OrgMemberRow{
		ID:             member.ID,
		AccountID:      member.AccountID,
		OrganizationID: member.OrganizationID,
		Head:           member.Head,
	})
	if err != nil {
		return models.OrgMember{}, err
	}

	return row.ToModel(), nil
}

func (r *Repository) DeleteOrgMember(ctx context.Context, memberID uuid.UUID) error {
	return r.OrgMembersQ.New().
		FilterByID(memberID).
		Delete(ctx)
}

func (r *Repository) GetOrgMemberByAccountID(
	ctx context.Context,
	organizationID, accountID uuid.UUID,
) (models.OrgMember, error) {
	row, err := r.OrgMembersQ.New().
		FilterByOrganizationID(organizationID).
		FilterByAccountID(accountID).
		Get(ctx)
	if err != nil {
		return models.OrgMember{}, err
	}

	return row.ToModel(), nil
}
