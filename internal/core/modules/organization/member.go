package organization

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/models"
)

func (m *Module) CreateOrgMember(
	ctx context.Context,
	member models.OrgMember,
) error {
	return m.repo.CreateOrgMember(ctx, member)
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
