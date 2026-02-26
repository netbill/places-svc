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
	OrganizationID uuid.UUID `json:"organization_id"`
	ClassID        uuid.UUID `json:"class_id"`
	Point          orb.Point `json:"point"`
	Address        string    `json:"address"`
	Name           string    `json:"name"`

	Description *string `json:"description"`
	Website     *string `json:"website"`
	Phone       *string `json:"phone"`
}

func (m *Module) Create(
	ctx context.Context,
	actor models.AccountActor,
	params CreateParams,
) (place models.Place, err error) {
	member, err := m.repo.GetOrgMemberByAccountID(ctx, actor, params.OrganizationID)
	if err != nil {
		return models.Place{}, err
	}
	if !member.Head {
		return models.Place{}, errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("only organization head can create places"),
		)
	}

	org, err := m.repo.GetOrganization(ctx, place.OrganizationID)
	if err != nil {
		return models.Place{}, err
	}
	if org.Status == models.OrganizationStatusSuspended {
		return models.Place{}, errx.ErrorOrganizationIsSuspended.Raise(
			fmt.Errorf("organization %s is suspended", place.OrganizationID),
		)
	}

	if !m.territory.ContainsLatLng(params.Point[1], params.Point[0]) {
		return models.Place{}, errx.ErrorPlaceOutOfTerritory.Raise(
			fmt.Errorf("place point %v is out of allowed territory", params.Point),
		)
	}

	class, err := m.repo.GetPlaceClass(ctx, params.ClassID)
	if err != nil {
		return models.Place{}, err
	}
	if class.DeprecatedAt != nil {
		return models.Place{}, errx.ErrorPlaceClassIsDeprecated.Raise(
			fmt.Errorf("place class %s is deprecated", params.ClassID),
		)
	}

	err = m.repo.Transaction(ctx, func(txCtx context.Context) error {
		place, err = m.repo.CreatePlace(ctx, params)
		if err != nil {
			return err
		}

		err = m.messenger.PublishCreatePlace(txCtx, place)
		if err != nil {
			return err
		}

		return nil
	})

	return place, nil
}
