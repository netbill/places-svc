package pclass

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
)

func (m *Module) Replace(ctx context.Context, oldClassID, newClassID uuid.UUID) error {
	if oldClassID == newClassID {
		return errx.ErrorInvalidInput.Raise(
			fmt.Errorf("old class ID and new class ID are the same: %s", oldClassID),
		)
	}

	_, err := m.Get(ctx, oldClassID)
	if err != nil {
		return err
	}

	_, err = m.Get(ctx, newClassID)
	if err != nil {
		return err
	}

	return m.repo.Transaction(ctx, func(txCtx context.Context) error {
		err = m.repo.ReplacePlacesClassID(txCtx, oldClassID, newClassID)
		if err != nil {
			return err
		}

		err = m.messanger.PublishPlacesClassReplaced(txCtx, oldClassID, newClassID)
		if err != nil {
			return err
		}

		return nil
	})
}
