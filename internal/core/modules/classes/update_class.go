package classes

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/models"
)

type UpdateParams struct {
	Code        *string `json:"code"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Icon        *string `json:"icon,omitempty"`
}

func (s Service) UpdateClassParams(
	ctx context.Context,
	ID uuid.UUID,
	updateParams UpdateParams,
) (class models.PlaceClass, err error) {
	if err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		class, err = s.repo.UpdatePlaceClass(ctx, ID, updateParams)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to update class %s: %w", ID.String(), err),
			)
		}

		if err = s.messanger.PublishPlaceClassUpdated(ctx, class); err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to publish class %s updated event: %w", ID.String(), err),
			)
		}

		return nil
	}); err != nil {
		return models.PlaceClass{}, err
	}

	return class, nil
}
