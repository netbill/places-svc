package pclass

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/models"
)

type CreateParams struct {
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
	Code        string     `json:"code"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Icon        *string    `json:"icon,omitempty"`
}

func (s Service) CreatePlaceClass(ctx context.Context, params CreateParams) (class models.PlaceClass, err error) {
	if params.ParentID != nil {
		parent, err := s.repo.GetPlaceClass(ctx, *params.ParentID)
		if err != nil {
			return models.PlaceClass{}, errx.ErrorInternal.Raise(
				fmt.Errorf("parent place class not found: %w", err),
			)
		}
		if parent.IsNil() {
			return models.PlaceClass{}, errx.ErrorPlaceClassNotFound.Raise(
				fmt.Errorf("parent place class not found"),
			)
		}

		cycle, err := s.repo.CheckParentCycle(ctx, uuid.Nil, *params.ParentID)
		if err != nil {
			return models.PlaceClass{}, errx.ErrorInternal.Raise(
				fmt.Errorf("parent cycle check failed: %w", err),
			)
		}
		if cycle {
			return models.PlaceClass{}, errx.ErrorPlaceClassParentCycle.Raise(
				fmt.Errorf("parent cycle detected"),
			)
		}
	}

	codeIsUsed, err := s.repo.PlaceClassExistsByCode(ctx, params.Code)
	if err != nil {
		return models.PlaceClass{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to check place class code uniqueness: %w", err),
		)
	}
	if codeIsUsed {
		return models.PlaceClass{}, errx.ErrorPlaceClassCodeExists.Raise(
			fmt.Errorf("place class code already in use"),
		)
	}

	if err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		class, err = s.repo.CreatePlaceClass(ctx, params)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to create place class: %w", err),
			)
		}

		err = s.messanger.PublishPlaceClassCreated(ctx, class)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to publish place class created event: %w", err),
			)
		}

		return nil
	}); err != nil {
		return models.PlaceClass{}, err
	}

	return class, nil
}
