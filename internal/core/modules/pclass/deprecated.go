package pclass

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/models"
)

func (m *Module) Deprecate(
	ctx context.Context,
	classID uuid.UUID,
	value bool,
) (models.PlaceClass, error) {
	class, err := m.repo.GetPlaceClass(ctx, classID)
	if err != nil {
		return models.PlaceClass{}, err
	}

	if class.DeprecatedAt != nil && value || class.DeprecatedAt == nil && !value {
		return class, nil
	}

	if err = m.repo.Transaction(ctx, func(ctx context.Context) error {
		class, err = m.repo.DeprecatedPlaceClass(ctx, classID, value)
		if err != nil {
			return err
		}

		return m.messenger.PublishPlaceClassUpdated(ctx, class)
	}); err != nil {
		return models.PlaceClass{}, err
	}

	return class, nil
}
