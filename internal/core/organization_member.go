package core

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
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

func (m *OrgModule) CreateMember(
	ctx context.Context,
	params CreateMemberParams,
) error {
	exists, err := m.member.ExistsByID(ctx, params.ID)
	if err != nil {
		return err
	}
	if exists {
		return errx.ErrorOrgMemberAlreadyExists.Raise(
			fmt.Errorf("organization member with id %s already exists", params.ID),
		)
	}

	bury, err := m.tombstone.OrgMemberIsBuried(ctx, params.ID)
	if err != nil {
		return err
	}
	if bury {
		return errx.ErrorOrgMemberDeleted.Raise(
			fmt.Errorf("organization member with id %s is already deleted", params.ID),
		)
	}

	bury, err = m.tombstone.OrganizationIsBuried(ctx, params.ID)
	if err != nil {
		return err
	}
	if bury {
		return errx.ErrorOrganizationDeleted.Raise(
			fmt.Errorf("organization with id %s is already deleted", params.OrganizationID),
		)
	}

	return m.member.Create(ctx, params)
}

type UpdateMemberParams struct {
	Label    *string `json:"label,omitempty"`
	Position *string `json:"position,omitempty"`

	Version   int32     `json:"version"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (m *OrgModule) UpdateMember(
	ctx context.Context,
	memberID uuid.UUID,
	params UpdateMemberParams,
) error {
	bury, err := m.tombstone.OrgMemberIsBuried(ctx, memberID)
	if err != nil {
		return err
	}
	if bury {
		return errx.ErrorOrgMemberDeleted.Raise(
			fmt.Errorf("organization member with id %s is already deleted", memberID),
		)
	}

	member, err := m.member.GetByID(ctx, memberID)
	if err != nil {
		return err
	}
	if params.Version <= member.Version {
		return nil
	}

	return m.member.Update(ctx, memberID, params)
}

func (m *OrgModule) DeleteMember(
	ctx context.Context,
	memberID uuid.UUID,
) error {
	bury, err := m.tombstone.OrgMemberIsBuried(ctx, memberID)
	if err != nil {
		return err
	}
	if bury {
		return errx.ErrorOrgMemberDeleted.Raise(
			fmt.Errorf("organization member with id %s is already deleted", memberID),
		)
	}

	return m.tx.Transaction(ctx, func(ctx context.Context) error {
		if err := m.tombstone.BuryOrgMember(ctx, memberID); err != nil {
			return err
		}

		return m.member.Delete(ctx, memberID)
	})
}
