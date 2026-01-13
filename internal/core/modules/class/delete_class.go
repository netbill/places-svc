package class

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
)

func (s Service) DeleteClass(ctx context.Context, classID uuid.UUID) error {
	exist, err := s.repo.CheckClassHasChildren(ctx, classID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to check if class exists: %w", err),
		)
	}
	if exist {
		return errx.ErrorClassHaveChildren.Raise(
			fmt.Errorf("cannot delete class %s because he haw child", classID),
		)
	}

	exist, err = s.repo.CheckPlaceExistForClass(ctx, classID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to check if class exists: %w", err),
		)
	}
	if !exist {
		return errx.ErrorPlacesExitsWithThisClass.Raise(
			fmt.Errorf("cannot delete class %s when place exist with thi class", classID),
		)
	}

	if err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		if err = s.repo.DeletePlaceClass(ctx, classID); err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to delete class %s: %w", classID, err),
			)
		}

		if err = s.messanger.PublishPlaceClassDeleted(ctx, classID); err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to delete class %s: %w", classID, err),
			)
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
