package pclass

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/models"
)

type CreateParams struct {
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Icon        *string    `json:"icon,omitempty"`
}

func (m *Module) Create(ctx context.Context, params CreateParams) (class models.PlaceClass, err error) {
	if params.ParentID != nil {
		if _, err = m.repo.GetPlaceClass(ctx, *params.ParentID); err != nil {
			return models.PlaceClass{}, err
		}
	}

	if err = m.repo.Transaction(ctx, func(ctx context.Context) error {
		class, err = m.repo.CreatePlaceClass(ctx, params)
		if err != nil {
			return err
		}

		err = m.messenger.PublishPlaceClassCreated(ctx, class)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return models.PlaceClass{}, err
	}

	return class, nil
}
