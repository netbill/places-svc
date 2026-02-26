package repository

import (
	"context"
)

type Repository struct {
	Transactioner
	PlacesQ        PlacesQ
	PlaceClassesQ  PlaceClassesQ
	OrganizationsQ OrganizationsQ
	OrgMembersQ    OrgMembersQ
}

type Transactioner interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
