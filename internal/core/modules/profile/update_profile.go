package profile

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/models"
)

type UpdateProfileParams struct {
	Official  bool
	Pseudonym *string
}

func (s Service) UpdateProfile(ctx context.Context, ID uuid.UUID, params UpdateProfileParams) (models.Profile, error) {
	res, err := s.repo.UpdateProfile(ctx, ID, params)
	if err != nil {
		return models.Profile{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to upsert profile: %w", err),
		)
	}

	return res, nil
}
