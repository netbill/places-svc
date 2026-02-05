package repository

import (
	"context"
)

type Repository struct {
	Transactioner
	OrganizationsQ          OrganizationsQ
	OrgMembersQ             OrgMembersQ
	OrgMemberRolesQ         OrgMemberRolesQ
	OrgRolesQ               OrgRolesQ
	OrgRolePermissionLinksQ OrgRolePermissionLinksQ
	PlacesQ                 PlacesQ
	PlaceClassesQ           PlaceClassesQ
}

type Transactioner interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
