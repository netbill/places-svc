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
	Name        string     `json:"name"`
	Description string     `json:"description"`
	IconKey     *string    `json:"icon_key,omitempty"`
}

func (m *Module) Update(
	ctx context.Context,
	classID uuid.UUID,
	params UpdateParams,
) (class models.PlaceClass, err error) {
	class, err = m.repo.GetPlaceClass(ctx, classID)
	if err != nil {
		return models.PlaceClass{}, err
	}

	if params.ParentID != nil {
		exist, err := m.repo.CheckParentCycle(ctx, class.ID, *class.ParentID)
		if err != nil {
			return models.PlaceClass{}, err
		}
		if exist {
			return models.PlaceClass{}, errx.ErrorPlaceClassParentCycle.Raise(
				fmt.Errorf("setting parent %s for class %s would create a cycle", *class.ParentID, class.ID),
			)
		}
	}

	icon, err := m.updatePlaceClassIcon(ctx, class, params)
	if err != nil {
		return models.PlaceClass{}, err
	}
	params.IconKey = icon

	if err = m.repo.Transaction(ctx, func(ctx context.Context) error {
		class, err = m.repo.UpdatePlaceClass(ctx, classID, params)
		if err != nil {
			return err
		}

		if err = m.messenger.PublishPlaceClassUpdated(ctx, class); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return models.PlaceClass{}, err
	}

	return class, nil
}
