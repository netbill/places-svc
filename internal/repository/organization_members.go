package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/models"
)

type OrgMemberRow struct {
	ID             uuid.UUID `db:"id"`
	AccountID      uuid.UUID `db:"account_id"`
	OrganizationID uuid.UUID `db:"organization_id"`
	Head           bool      `db:"head"`
	Position       *string   `db:"position,omitempty"`
	Label          *string   `db:"label,omitempty"`

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
		Position:       r.Position,
		Label:          r.Label,
		CreatedAt:      r.SourceCreatedAt,
		UpdatedAt:      r.SourceUpdatedAt,
	}
}

type OrgMembersQ interface {
	New() OrgMembersQ
	Insert(ctx context.Context, input OrgMemberRow) (OrgMemberRow, error)

	Get(ctx context.Context) (OrgMemberRow, error)
	Select(ctx context.Context) ([]OrgMemberRow, error)

	UpdateMany(ctx context.Context) (int64, error)
	UpdateOne(ctx context.Context) (OrgMemberRow, error)

	UpdateHead(head bool) OrgMembersQ
	UpdatePosition(position *string) OrgMembersQ
	UpdateLabel(label *string) OrgMembersQ
	UpdateSourceUpdatedAt(v time.Time) OrgMembersQ

	FilterByID(id ...uuid.UUID) OrgMembersQ
	FilterByAccountID(accountID ...uuid.UUID) OrgMembersQ
	FilterByOrganizationID(organizationID ...uuid.UUID) OrgMembersQ
	FilterByHead(head bool) OrgMembersQ

	Delete(ctx context.Context) error
}
