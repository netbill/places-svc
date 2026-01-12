package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/models"
	"github.com/netbill/places-svc/internal/repository/pgdb"
)

func (s Service) UpsertOrganization(ctx context.Context, params models.Organization) error {
	_, err := s.organizationsQ(ctx).Upsert(ctx, pgdb.Organization{
		ID:        params.ID,
		Status:    params.Status,
		Verified:  params.Verified,
		Name:      params.Name,
		Icon:      params.Icon,
		CreatedAt: params.CreatedAt,
		UpdatedAt: params.UpdatedAt,
	})

	return err
}

func (s Service) DeleteOrganization(ctx context.Context, ID uuid.UUID) error {
	return s.organizationsQ(ctx).FilterByID(ID).Delete(ctx)
}

func (s Service) UpsertOrgMember(ctx context.Context, member models.Member) error {
	err := s.orgMembersQ(ctx).Upsert(ctx, pgdb.OrgMember{
		ID:             member.ID,
		OrganizationID: member.OrganizationID,
		AccountID:      member.AccountID,
		CreatedAt:      member.CreatedAt,
		UpdatedAt:      member.UpdatedAt,
	})

	return err
}

func (s Service) DeleteOrgMember(ctx context.Context, ID uuid.UUID) error {
	return s.orgMembersQ(ctx).FilterByID(ID).Delete(ctx)
}

func (s Service) UpsertOrgRole(ctx context.Context, params models.OrgRole) error {
	err := s.orgRolesQ(ctx).Upsert(ctx, pgdb.OrgRole{
		ID:             params.ID,
		OrganizationID: params.OrganizationID,
		Rank:           params.Rank,
		Head:           params.Head,
		CreatedAt:      params.CreatedAt,
		UpdatedAt:      params.UpdatedAt,
	})

	return err
}

func (s Service) DeleteOrgRole(ctx context.Context, ID uuid.UUID) error {
	return s.orgRolesQ(ctx).FilterByID(ID).Delete(ctx)
}

func (s Service) UpdateOrgRolesRanks(
	ctx context.Context,
	organizationID uuid.UUID,
	order map[uuid.UUID]uint,
) error {
	_, err := s.orgRolesQ(ctx).UpdateRolesRanks(ctx, organizationID, order)
	if err != nil {
		return err
	}

	return nil
}

func (s Service) UpdateOrgRolePermissions(
	ctx context.Context,
	roleID uuid.UUID,
	permissions map[string]bool,
) error {
	deletePermissions := make([]string, 0)
	addPermissions := make([]string, 0)

	for perm, toSet := range permissions {
		if toSet {
			addPermissions = append(addPermissions, perm)
		} else {
			deletePermissions = append(deletePermissions, perm)
		}
	}

	if len(deletePermissions) > 0 {
		if err := s.orgRolePermissionLinksQ(ctx).
			FilterByRoleID(roleID).
			FilterByPermissionCode(deletePermissions...).
			Delete(ctx); err != nil {
			return err
		}
	}

	if len(addPermissions) > 0 {
		p, err := s.orgRolePermissionsQ(ctx).FilterByCode(addPermissions...).Select(ctx)
		if err != nil {
			return err
		}

		existingPermissionsMap := make([]pgdb.OrganizationRolePermissionLink, len(p))
		for i, perm := range p {
			existingPermissionsMap[i] = pgdb.OrganizationRolePermissionLink{
				RoleID:         roleID,
				PermissionCode: perm.Code,
			}
		}
		if err = s.orgRolePermissionLinksQ(ctx).Insert(ctx, existingPermissionsMap...); err != nil {
			return err
		}
	}

	return nil
}

func (s Service) CheckMemberHavePermission(
	ctx context.Context,
	memberID uuid.UUID,
	permissionCode string,
) (bool, error) {
	have, err := s.orgMembersQ(ctx).
		FilterByID(memberID).
		FilterByPermissionCode(permissionCode).
		Exists(ctx)
	if err != nil {
		return false, fmt.Errorf("checking member have permission: %w", err)
	}

	return have, nil
}

func (s Service) RemoveOrgMemberRole(ctx context.Context, memberID, roleID uuid.UUID) error {
	return s.orgMemberRolesQ(ctx).
		FilterByMemberID(memberID).
		FilterByRoleID(roleID).
		Delete(ctx)
}

func (s Service) AddOrgMemberRole(ctx context.Context, memberID, roleID uuid.UUID) error {
	_, err := s.orgMemberRolesQ(ctx).
		Insert(ctx, pgdb.OrgMemberRole{
			MemberID: memberID,
			RoleID:   roleID,
		})

	return err
}
