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
	Website     *string `json:"website"`
	Phone       *string `json:"phone"`
}

func (m *Module) CreatePlace(
	ctx context.Context,
	initiator models.InitiatorData,
	params CreateParams,
) (place models.Place, err error) {
	if params.OrganizationID != nil {
		err = m.chekPermissionForManagePlace(ctx, initiator, *params.OrganizationID)
		if err != nil {
			return models.Place{}, err
		}
	}

	if !m.territory.ContainsLatLng(params.Point[1], params.Point[0]) {
		return models.Place{}, errx.ErrorPlaceOutOfTerritory.Raise(
			fmt.Errorf("place point %v is out of allowed territory", params.Point),
		)
	}

	classExists, err := m.repo.CheckPlaceClassExists(ctx, params.ClassID)
	if err != nil {
		return models.Place{}, err
	}
	if !classExists {
		return models.Place{}, errx.ErrorPlaceClassNotFound.Raise(
			fmt.Errorf("place class %v not found", params.ClassID),
		)
	}

	err = m.repo.Transaction(ctx, func(txCtx context.Context) error {
		place, err = m.repo.CreatePlace(ctx, params)
		if err != nil {
			return err
		}

		err = m.messanger.PublishCreatePlace(txCtx, place)
		if err != nil {
			return err
		}

		return nil
	})

	return place, nil
}
