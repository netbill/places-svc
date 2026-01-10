package organization

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
)

func (s Service) AddMemberRole(
	ctx context.Context,
	memberID, roleID uuid.UUID,
) error {
	if err := s.repo.AddMemberRole(ctx, memberID, roleID); err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("role Service MemberAddRoleByUser: repo AddMemberRole: %w", err),
		)
	}

	return nil
}

func (s Service) RemoveMemberRole(
	ctx context.Context,
	memberID, roleID uuid.UUID,
) error {
	if err := s.repo.RemoveMemberRole(ctx, memberID, roleID); err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("role Service MemberRemoveRoleByUser: repo RemoveMemberRole: %w", err),
		)
	}

	return nil
}
