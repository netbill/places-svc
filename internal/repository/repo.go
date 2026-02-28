package repository

import (
	"context"
)

type Repository struct {
	PlacesSql        PlacesQ
	PlaceClassesSql  PlaceClassesQ
	OrganizationsSql OrganizationsQ
	OrgMembersSql    OrgMembersQ
	TombstonesSql
	TransactionSql
}

type TransactionSql interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
