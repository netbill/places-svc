package organization

import (
	"context"
	"time"

	"github.com/google/uuid"
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
	return m.repo.UpdateOrgMember(ctx, memberID, params)
}

func (m *Module) DeleteOrgMember(
	ctx context.Context,
	memberID uuid.UUID,
) error {
	return m.repo.DeleteOrgMember(ctx, memberID)
}
