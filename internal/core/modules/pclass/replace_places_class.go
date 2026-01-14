package pclass

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
)

func (s Service) ReplacePlacesClass(ctx context.Context, oldClassID, newClassID uuid.UUID) error {
	if oldClassID == newClassID {
		return errx.ErrorInvalidInput.Raise(
			fmt.Errorf("old class ID and new class ID are the same: %s", oldClassID),
		)
	}

	_, err := s.GetPlaceClass(ctx, oldClassID)
	if err != nil {
		return err
	}

	_, err = s.GetPlaceClass(ctx, newClassID)
	if err != nil {
		return err
	}

	return s.repo.Transaction(ctx, func(txCtx context.Context) error {
		err = s.repo.ReplacePlacesClassID(txCtx, oldClassID, newClassID)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to replace places class ID from %s to %s: %w", oldClassID, newClassID, err),
			)
		}

		err = s.messanger.PublishPlacesClassReplaced(txCtx, oldClassID, newClassID)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf(
					"failed to publish places class replaced event for old class ID %s and new class ID %s: %w",
					oldClassID, newClassID, err,
				),
			)
		}

		return nil
	})
}
