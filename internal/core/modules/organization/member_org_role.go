package organization

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
)

func (s Service) AddOrgMemberRole(
	ctx context.Context,
	memberID, roleID uuid.UUID,
) error {
	if err := s.repo.AddOrgMemberRole(ctx, memberID, roleID); err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("role Service MemberAddRoleByUser: repo AddOrgMemberRole: %w", err),
		)
	}

	return nil
}

func (s Service) RemoveOrgMemberRole(
	ctx context.Context,
	memberID, roleID uuid.UUID,
) error {
	if err := s.repo.RemoveOrgMemberRole(ctx, memberID, roleID); err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("role Service MemberRemoveRoleByUser: repo RemoveOrgMemberRole: %w", err),
		)
	}

	return nil
}
