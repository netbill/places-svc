package organization

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

func (m *Module) CreateOrgMember(
	ctx context.Context,
	params CreateMemberParams,
) error {
	exists, err := m.repo.ExistsOrgMember(ctx, params.ID)
	if err != nil {
		return err
	}
	if exists {
		return errx.ErrorOrgMemberAlreadyExists.Raise(
			fmt.Errorf("organization member with id %s already exists", params.ID),
		)
	}

	bury, err := m.repo.OrgMemberIsBuried(ctx, params.ID)
	if err != nil {
		return err
	}
	if bury {
		return errx.ErrorOrgMemberDeleted.Raise(
			fmt.Errorf("organization member with id %s is already deleted", params.ID),
		)
	}

	bury, err = m.repo.OrganizationIsBuried(ctx, params.ID)
	if err != nil {
		return err
	}
	if bury {
		return errx.ErrorOrganizationDeleted.Raise(
			fmt.Errorf("organization with id %s is already deleted", params.OrganizationID),
		)
	}

	return m.repo.CreateOrgMember(ctx, params)
}

type UpdateMemberParams struct {
	Label    *string `json:"label,omitempty"`
	Position *string `json:"position,omitempty"`

	Version   int32     `json:"version"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (m *Module) UpdateOrgMember(
	ctx context.Context,
	memberID uuid.UUID,
	params UpdateMemberParams,
) error {
	bury, err := m.repo.OrgMemberIsBuried(ctx, memberID)
	if err != nil {
		return err
	}
	if bury {
		return errx.ErrorOrgMemberDeleted.Raise(
			fmt.Errorf("organization member with id %s is already deleted", memberID),
		)
	}

	return m.repo.UpdateOrgMember(ctx, memberID, params)
}

func (m *Module) DeleteOrgMember(
	ctx context.Context,
	memberID uuid.UUID,
) error {
	bury, err := m.repo.OrgMemberIsBuried(ctx, memberID)
	if err != nil {
		return err
	}
	if bury {
		return errx.ErrorOrgMemberDeleted.Raise(
			fmt.Errorf("organization member with id %s is already deleted", memberID),
		)
	}

	return m.repo.DeleteOrgMember(ctx, memberID)
}
