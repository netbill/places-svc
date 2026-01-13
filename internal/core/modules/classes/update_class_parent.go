package classes

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/models"
)

func (s Service) UpdateClassParent(
	ctx context.Context,
	classID uuid.UUID,
	parentID *uuid.UUID,
) (class models.PlaceClass, err error) {
	if err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		if parentID != nil {
			exist, err := s.repo.CheckParentCycle(ctx, classID, *parentID)
			if err != nil {
				return errx.ErrorInternal.Raise(
					fmt.Errorf("failed to check parent cycle for class %s and parent %s: %w", classID, *parentID, err),
				)
			}
			if exist {
				return errx.ErrorPlaceClassParentCycle.Raise(
					fmt.Errorf("setting parent %s for class %s would create a cycle", *parentID, classID),
				)
			}
		}

		class, err = s.repo.UpdatePlaceClassParent(ctx, classID, parentID)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to update class %s parent: %w", classID, err),
			)
		}

		if err = s.messanger.PublishPlaceClassParentUpdated(ctx, class); err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to publish class %s parent updated event: %w", classID, err),
			)
		}

		return nil
	}); err != nil {
		return models.PlaceClass{}, err
	}

	return class, nil
}
