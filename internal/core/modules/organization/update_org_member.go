package organization

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/models"
)

type UpdateMemberParams struct {
	Position  *string   `json:"position,omitempty"`
	Label     *string   `json:"label,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (m *Module) UpdateOrgMember(
	ctx context.Context,
	memberID uuid.UUID,
	params UpdateMemberParams,
) (models.OrgMember, error) {
	return m.repo.UpdateOrgMember(ctx, memberID, params)
}
