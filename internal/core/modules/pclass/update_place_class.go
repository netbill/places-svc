package pclass

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/models"
)

type UpdateParams struct {
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
	Code        *string    `json:"code"`
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	Icon        *string    `json:"icon,omitempty"`
}

func (s Service) UpdatePlaceClass(
	ctx context.Context,
	ID uuid.UUID,
	params UpdateParams,
) (class models.PlaceClass, err error) {
	class, err = s.repo.GetPlaceClass(ctx, ID)
	if err != nil {
		return models.PlaceClass{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get class %s: %w", ID.String(), err),
		)
	}
	if class.IsNil() {
		return models.PlaceClass{}, errx.ErrorPlaceClassNotFound.Raise(
			fmt.Errorf("place class %s not found", ID.String()),
		)
	}

	if params.ParentID != nil {
		if *params.ParentID != uuid.Nil {
			exist, err := s.repo.CheckParentCycle(ctx, class.ID, *class.ParentID)
			if err != nil {
				return models.PlaceClass{}, errx.ErrorInternal.Raise(
					fmt.Errorf("failed to check parent cycle for class %s and parent %s: %w", class.ID, *class.ParentID, err),
				)
			}
			if exist {
				return models.PlaceClass{}, errx.ErrorPlaceClassParentCycle.Raise(
					fmt.Errorf("setting parent %s for class %s would create a cycle", *class.ParentID, class.ID),
				)
			}
		}
	}

	if params.Code != nil {
		codeIsUsed, err := s.repo.PlaceClassExistsByCode(ctx, *params.Code)
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
	}

	if err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		class, err = s.repo.UpdatePlaceClass(ctx, ID, params)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to update class %s: %w", ID.String(), err),
			)
		}

		if err = s.messanger.PublishPlaceClassUpdated(ctx, class); err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to publish class %s updated event: %w", ID.String(), err),
			)
		}

		return nil
	}); err != nil {
		return models.PlaceClass{}, err
	}

	return class, nil
}
