package place

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/models"
	"github.com/paulmach/orb"
)

type CreateParams struct {
	OrganizationID *uuid.UUID `json:"organization_id"`
	ClassID        uuid.UUID  `json:"class_id"`
	Point          orb.Point  `json:"point"`
	Address        string     `json:"address"`
	Name           string     `json:"name"`

	Description *string `json:"description"`
	Icon        *string `json:"icon"`
	Banner      *string `json:"banner"`
	Website     *string `json:"website"`
	Phone       *string `json:"phone"`
}

func (s Service) CreatePlace(
	ctx context.Context,
	initiatorID uuid.UUID,
	params CreateParams,
) (place models.Place, err error) {
	if params.OrganizationID != nil {
		err = s.chekPermissionForManagePlace(ctx, initiatorID, *params.OrganizationID)
		if err != nil {
			return models.Place{}, err
		}
	}

	if !s.territory.ContainsLatLng(params.Point[1], params.Point[0]) {
		return models.Place{}, errx.ErrorPlaceOutOfTerritory.Raise(
			fmt.Errorf("place point %v is out of allowed territory", params.Point),
		)
	}

	classExists, err := s.repo.CheckPlaceClassExists(ctx, params.ClassID)
	if err != nil {
		return models.Place{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to check place class exists: %w", err),
		)
	}
	if !classExists {
		return models.Place{}, errx.ErrorPlaceClassNotFound.Raise(
			fmt.Errorf("place class %v not found", params.ClassID),
		)
	}

	err = s.repo.Transaction(ctx, func(txCtx context.Context) error {
		place, err = s.repo.CreatePlace(ctx, params)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to create place: %w", err),
			)
		}

		err = s.messanger.PublishCreatePlace(txCtx, place)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to publish create place message: %w", err),
			)
		}

		return nil
	})

	return place, nil
}
