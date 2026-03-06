package organization

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/errx"
	"github.com/netbill/places-svc/internal/models"
)

type CreateParams struct {
	ID        uuid.UUID `json:"id"`
	Status    string    `json:"status"`
	Name      string    `json:"name"`
	IconKey   *string   `json:"icon_key,omitempty"`
	BannerKey *string   `json:"banner_key,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}

func (s *Service) Create(
	ctx context.Context,
	org CreateParams,
) error {
	exists, err := s.org.Exists(ctx, org.ID)
	if err != nil {
		return err
	}
	if exists {
		return errx.ErrorOrganizationAlreadyExists.Raise(
			fmt.Errorf("orgRepo with id %s already exists", org.ID),
		)
	}

	bury, err := s.tombstone.OrganizationIsBuried(ctx, org.ID)
	if err != nil {
		return err
	}
	if bury {
		return errx.ErrorOrganizationDeleted.Raise(
			fmt.Errorf("orgRepo with id %s is already deleted", org.ID),
		)
	}

	return s.org.Create(ctx, org)
}

func (s *Service) Get(
	ctx context.Context,
	orgID uuid.UUID,
) (models.Organization, error) {
	res, err := s.org.Get(ctx, orgID)
	if errors.Is(err, errx.ErrorOrganizationNotExists) {
		buried, err := s.tombstone.OrganizationIsBuried(ctx, orgID)
		if err != nil {
			return models.Organization{}, err
		}
		if buried {
			return models.Organization{}, errx.ErrorOrganizationDeleted.Raise(
				fmt.Errorf("orgRepo with id %s is deleted", orgID),
			)
		}
	}
	if err != nil {
		return models.Organization{}, err
	}

	return res, nil
}

func (s *Service) GetByIDs(
	ctx context.Context,
	ids []uuid.UUID,
) ([]models.Organization, error) {
	return s.org.GetListByIDs(ctx, ids)
}

type UpdateParams struct {
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	IconKey   *string   `json:"icon_key"`
	BannerKey *string   `json:"banner_key"`
	Version   int32     `json:"version"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *Service) Update(
	ctx context.Context,
	orgID uuid.UUID,
	params UpdateParams,
) error {
	org, err := s.Get(ctx, orgID)
	if err != nil {
		return err
	}

	if params.Version <= org.Version {
		return nil
	}

	return s.org.Update(ctx, orgID, params)
}

func (s *Service) Delete(
	ctx context.Context,
	organizationID uuid.UUID,
) error {
	buried, err := s.tombstone.OrganizationIsBuried(ctx, organizationID)
	if err != nil {
		return err
	}
	if buried {
		return errx.ErrorOrganizationDeleted.Raise(
			fmt.Errorf("orgRepo with id %s is already deleted", organizationID),
		)
	}

	return s.tx.Transaction(ctx, func(ctx context.Context) error {
		if err := s.tombstone.BuryOrganization(ctx, organizationID); err != nil {
			return err
		}

		return s.org.Delete(ctx, organizationID)
	})
}
