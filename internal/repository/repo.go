package repository

import (
	"context"
	"database/sql"

	"github.com/netbill/pgx"
	replicaspg "github.com/netbill/replicas/pgdb"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r Repository) OrganizationsQ(ctx context.Context) replicaspg.OrganizationsQ {
	return replicaspg.NewOrganizationsQ(pgx.Exec(r.db, ctx))
}

func (r Repository) membersQ(ctx context.Context) replicaspg.OrgMembersQ {
	return replicaspg.NewOrgMembersQ(pgx.Exec(r.db, ctx))
}

func (r Repository) memberRolesQ(ctx context.Context) replicaspg.OrgMemberRolesQ {
	return replicaspg.NewOrgMemberRolesQ(pgx.Exec(r.db, ctx))
}

func (r Repository) rolesQ(ctx context.Context) replicaspg.OrgRolesQ {
	return replicaspg.NewOrgRolesQ(pgx.Exec(r.db, ctx))
}

func (r Repository) rolePermissionLinksQ(ctx context.Context) replicaspg.OrgRolePermissionLinksQ {
	return replicaspg.NewOrgRolePermissionsQ(pgx.Exec(r.db, ctx))
}

func (r Repository) rolePermissionsQ(ctx context.Context) replicaspg.OrgRolePermissionsQ {
	return replicaspg.NewOrgPermissionsQ(pgx.Exec(r.db, ctx))
}

func (r Repository) invitesQ(ctx context.Context) replicaspg.OrgInvitesQ {
	return replicaspg.NewOrgInvitesQ(pgx.Exec(r.db, ctx))
}

func (r Repository) profilesQ(ctx context.Context) replicaspg.ProfilesQ {
	return replicaspg.NewProfilesQ(pgx.Exec(r.db, ctx))
}

func (r Repository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return pgx.Transaction(r.db, ctx, fn)
}
