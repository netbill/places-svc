package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/models"
)

type OrgMemberRoleRow struct {
	MemberID         uuid.UUID `db:"member_id"`
	RoleID           uuid.UUID `db:"role_id"`
	SourceCreatedAt  time.Time `db:"source_created_at"`
	ReplicaCreatedAt time.Time `db:"replica_created_at"`
}

func (r OrgMemberRoleRow) IsNil() bool {
	return r.MemberID == uuid.Nil && r.RoleID == uuid.Nil
}

func (r OrgMemberRoleRow) ToModel() models.OrgMemberRolesLink {
	return models.OrgMemberRolesLink{
		MemberID:  r.MemberID,
		RoleID:    r.RoleID,
		CreatedAt: r.SourceCreatedAt,
	}
}

type OrgMemberRolesQ interface {
	New() OrgMemberRolesQ
	Insert(ctx context.Context, input OrgMemberRoleRow) (OrgMemberRoleRow, error)

	Get(ctx context.Context) (OrgMemberRoleRow, error)
	Select(ctx context.Context) ([]OrgMemberRoleRow, error)

	UpdateMany(ctx context.Context) (int64, error)
	UpdateOne(ctx context.Context) (OrgMemberRoleRow, error)

	UpdateSourceCreatedAt(v time.Time) OrgMemberRolesQ

	FilterByMemberID(memberID ...uuid.UUID) OrgMemberRolesQ
	FilterByRoleID(roleID ...uuid.UUID) OrgMemberRolesQ

	Delete(ctx context.Context) error
}
