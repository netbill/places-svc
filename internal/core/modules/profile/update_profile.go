package profile

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type UpdateParams struct {
	Username  string    `json:"username"`
	Official  bool      `json:"official"`
	Pseudonym *string   `json:"pseudonym"`
	Avatar    *string   `json:"avatar"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (m *Module) UpdateProfile(ctx context.Context, accountID uuid.UUID, params UpdateParams) error {
	return m.repo.UpdateProfile(ctx, accountID, params)
}
