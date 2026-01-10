package organization

import (
	"context"
	"fmt"

	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/models"
)

func (s Service) CreateRole(
	ctx context.Context,
	params models.Role,
) error {
	if err := s.repo.CreateRole(ctx, params); err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to create role: %w", err),
		)
	}

	return nil
}
