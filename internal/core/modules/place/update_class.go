package place

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/models"
)

func (m *Module) UpdateClass(
	ctx context.Context,
	initiator models.Initiator,
	placeID uuid.UUID,
	classID uuid.UUID,
) (place models.Place, err error) {
	place, err = m.repo.GetPlaceByID(ctx, placeID)
	if err != nil {
		return models.Place{}, err
	}

	err = m.chekPermissionForUpdatePlace(ctx, initiator, placeID)
	if err != nil {
		return models.Place{}, err
	}

	exist, err := m.repo.CheckPlaceClassExists(ctx, classID)
	if err != nil {
		return models.Place{}, err
	}
	if !exist {
		return models.Place{}, errx.ErrorPlaceClassCodeExists.Raise(
			fmt.Errorf("place class %s does not exist", classID),
		)
	}

	err = m.repo.Transaction(ctx, func(txCtx context.Context) error {
		place, err = m.repo.UpdateClassForPlace(ctx, placeID, classID)
		if err != nil {
			return err
		}

		err = m.messenger.PublishUpdatePlaceVerified(ctx, place)
		if err != nil {
			return err
		}

		return nil
	})

	return place, nil
}
