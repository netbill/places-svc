package class

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/models"
)

type CreateParams struct {
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
	Code        string     `json:"code"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Icon        *string    `json:"icon,omitempty"`
}

func (s Service) CreateClass(ctx context.Context, params CreateParams) (class models.PlaceClass, err error) {
	if err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		class, err = s.repo.CreatePlaceClass(ctx, params)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to create place class: %w", err),
			)
		}

		err = s.messanger.PublishPlaceClassCreated(ctx, class)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to publish place class created event: %w", err),
			)
		}

		return nil
	}); err != nil {
		return models.PlaceClass{}, err
	}

	return class, nil
}
