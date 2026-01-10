package organization

import (
	"context"
	"fmt"

	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/models"
)

func (s Service) CreateMember(ctx context.Context, member models.Member) error {
	err := s.repo.CreateMember(ctx, member)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to create organization member %w", err),
		)
	}

	return nil
}
