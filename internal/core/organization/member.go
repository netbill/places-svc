package organization

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/errx"
)

type CreateMemberParams struct {
	ID             uuid.UUID `json:"id"`
	AccountID      uuid.UUID `json:"account_id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Head           bool      `json:"head"`
	Label          *string   `json:"label,omitempty"`
	Position       *string   `json:"position,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}

func (s *Service) CreateMember(
	ctx context.Context,
	params CreateMemberParams,
) error {
	exists, err := s.member.ExistsByID(ctx, params.ID)
	if err != nil {
		return err
	}
	if exists {
		return errx.ErrorOrgMemberAlreadyExists.Raise(
			fmt.Errorf("orgRepo memberRepo with id %s already exists", params.ID),
		)
	}

	bury, err := s.tombstone.OrgMemberIsBuried(ctx, params.ID)
	if err != nil {
		return err
	}
	if bury {
		return errx.ErrorOrgMemberDeleted.Raise(
			fmt.Errorf("orgRepo memberRepo with id %s is already deleted", params.ID),
		)
	}

	bury, err = s.tombstone.OrganizationIsBuried(ctx, params.ID)
	if err != nil {
		return err
	}
	if bury {
		return errx.ErrorOrganizationDeleted.Raise(
			fmt.Errorf("orgRepo with id %s is already deleted", params.OrganizationID),
		)
	}

	return s.member.Create(ctx, params)
}

type UpdateMemberParams struct {
	Label    *string `json:"label,omitempty"`
	Position *string `json:"position,omitempty"`

	Version   int32     `json:"version"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *Service) UpdateMember(
	ctx context.Context,
	memberID uuid.UUID,
	params UpdateMemberParams,
) error {
	bury, err := s.tombstone.OrgMemberIsBuried(ctx, memberID)
	if err != nil {
		return err
	}
	if bury {
		return errx.ErrorOrgMemberDeleted.Raise(
			fmt.Errorf("orgRepo memberRepo with id %s is already deleted", memberID),
		)
	}

	member, err := s.member.GetByID(ctx, memberID)
	if err != nil {
		return err
	}
	if params.Version <= member.Version {
		return nil
	}

	return s.member.Update(ctx, memberID, params)
}

func (s *Service) DeleteMember(
	ctx context.Context,
	memberID uuid.UUID,
) error {
	bury, err := s.tombstone.OrgMemberIsBuried(ctx, memberID)
	if err != nil {
		return err
	}
	if bury {
		return errx.ErrorOrgMemberDeleted.Raise(
			fmt.Errorf("orgRepo memberRepo with id %s is already deleted", memberID),
		)
	}

	return s.tx.Transaction(ctx, func(ctx context.Context) error {
		if err := s.tombstone.BuryOrgMember(ctx, memberID); err != nil {
			return err
		}

		return s.member.Delete(ctx, memberID)
	})
}
