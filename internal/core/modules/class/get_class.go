package class

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/pagi"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/models"
)

func (s Service) GetClass(ctx context.Context, id uuid.UUID) (models.PlaceClass, error) {
	class, err := s.repo.GetClass(ctx, id)
	if err != nil {
		return models.PlaceClass{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get class %s: %w", id, err),
		)
	}
	if class.IsNil() {
		return models.PlaceClass{}, errx.ErrorPlaceClassNotFound.Raise(
			fmt.Errorf("class with id %s not found", id),
		)
	}

	return class, nil
}

type FilterParams struct {
	Parent      *FilterClassParams
	Text        *string
	Name        *string
	Description *string
}

type FilterClassParams struct {
	ParentID uuid.UUID
	Parents  bool
	Children bool
}

func (s Service) GetClasses(ctx context.Context, params FilterParams, limit, offset uint) (pagi.Page[[]models.PlaceClass], error) {
	if limit == 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	classes, err := s.repo.GetClasses(ctx, params, limit, offset)
	if err != nil {
		return pagi.Page[[]models.PlaceClass]{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get classes: %w", err),
		)
	}

	return classes, nil
}
