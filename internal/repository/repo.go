package repository

import (
	"context"
	"database/sql"

	"github.com/netbill/pgx"
	"github.com/netbill/places-svc/internal/repository/pgdb"
)

type Service struct {
	db *sql.DB
}

func New(db *sql.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s Service) organizationsQ(ctx context.Context) pgdb.OrganizationsQ {
	return pgdb.NewOrganizationsQ(pgx.Exec(s.db, ctx))
}

func (s Service) orgMembersQ(ctx context.Context) pgdb.OrgMembersQ {
	return pgdb.NewOrgMembersQ(pgx.Exec(s.db, ctx))
}

func (s Service) orgMemberRolesQ(ctx context.Context) pgdb.OrgMemberRolesQ {
	return pgdb.NewOrgMemberRolesQ(pgx.Exec(s.db, ctx))
}

func (s Service) orgRolesQ(ctx context.Context) pgdb.OrgRolesQ {
	return pgdb.NewOrgRolesQ(pgx.Exec(s.db, ctx))
}

func (s Service) orgRolePermissionLinksQ(ctx context.Context) pgdb.OrgRolePermissionLinksQ {
	return pgdb.NewOrgRolePermissionLinksQ(pgx.Exec(s.db, ctx))
}

func (s Service) orgRolePermissionsQ(ctx context.Context) pgdb.OrgRolePermissionsQ {
	return pgdb.NewOrgRolePermissionsQ(pgx.Exec(s.db, ctx))
}

func (s Service) profilesQ(ctx context.Context) pgdb.ProfilesQ {
	return pgdb.NewProfilesQ(pgx.Exec(s.db, ctx))
}

func (s Service) placeQ(ctx context.Context) pgdb.PlacesQ {
	return pgdb.NewPlacesQ(pgx.Exec(s.db, ctx))
}

func (s Service) placeClassesQ(ctx context.Context) pgdb.PlaceClassesQ {
	return pgdb.NewPlaceClassesQ(pgx.Exec(s.db, ctx))
}

func (s Service) placePossibilitiesQ(ctx context.Context) pgdb.PlacePossibilityLinksQ {
	return pgdb.NewPlacePossibilitiesQ(pgx.Exec(s.db, ctx))
}

func (s Service) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return pgx.Transaction(s.db, ctx, fn)
}
