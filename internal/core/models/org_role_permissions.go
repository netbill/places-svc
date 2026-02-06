package models

type OrgRolePermissionCode string

const (
	RolePermissionPlaceCreate = "place.create"
	RolePermissionPlaceDelete = "place.delete"
	RolePermissionPlaceUpdate = "place.update"
)

var allRolePermissions = []string{
	RolePermissionPlaceCreate,
	RolePermissionPlaceDelete,
	RolePermissionPlaceUpdate,
}

func GetOrgRolePermissionLength() int {
	return len(allRolePermissions)
}

type OrgRolePermission struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

type OrgRolePermissionDetails struct {
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

type OrgRolePermissionAccess struct {
	PlaceCreate bool `json:"create.place"`
	PlaceDelete bool `json:"delete.place"`
	PlaceUpdate bool `json:"update.place"`
}

func (p OrgRolePermissionAccess) ToMap() map[string]bool {
	perms := make(map[string]bool)

	perms[RolePermissionPlaceCreate] = p.PlaceCreate
	perms[RolePermissionPlaceDelete] = p.PlaceDelete
	perms[RolePermissionPlaceUpdate] = p.PlaceUpdate

	return perms
}

type OrgRolePermissionDictWithDetails struct {
	PlaceCreate OrgRolePermissionDetails `json:"place.create"`
	PlaceDelete OrgRolePermissionDetails `json:"place.delete"`
	PlaceUpdate OrgRolePermissionDetails `json:"place.update"`
}

func (p OrgRolePermissionDictWithDetails) ToMap() map[string]OrgRolePermissionDetails {
	perms := make(map[string]OrgRolePermissionDetails)

	perms[RolePermissionPlaceCreate] = p.PlaceCreate
	perms[RolePermissionPlaceDelete] = p.PlaceDelete
	perms[RolePermissionPlaceUpdate] = p.PlaceUpdate

	return perms
}
