package pclass

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
)

func (m *Module) Delete(ctx context.Context, classID uuid.UUID) error {
	exist, err := m.repo.PlaceClassExists(ctx, classID)
	if err != nil {
		return err
	}
	if !exist {
		return errx.ErrorPlaceClassNotExists.Raise(
			fmt.Errorf("place class with id %s not found", classID.String()),
		)
	}

	exist, err = m.repo.CheckPlaceClassHasChildren(ctx, classID)
	if err != nil {
		return err
	}
	if exist {
		return errx.ErrorPlaceClassHaveChildren.Raise(
			fmt.Errorf("cannot delete class %s because he haw child", classID),
		)
	}

	exist, err = m.repo.CheckPlaceExistForClass(ctx, classID)
	if err != nil {
		return err
	}
	if !exist {
		return errx.ErrorPlacesExitsWithThisClass.Raise(
			fmt.Errorf("cannot delete class %s when place exist with thi class", classID),
		)
	}

	if err = m.repo.Transaction(ctx, func(ctx context.Context) error {
		if err = m.repo.DeletePlaceClass(ctx, classID); err != nil {
			return err
		}

		if err = m.messenger.PublishPlaceClassDeleted(ctx, classID); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
