package core

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/models"
)

type orgRepo interface {
	Create(
		ctx context.Context,
		params CreateOrgParams,
	) error

	Get(
		ctx context.Context,
		orgID uuid.UUID,
	) (models.Organization, error)
	GetListByIDs(
		ctx context.Context,
		ids []uuid.UUID,
	) ([]models.Organization, error)

	Exists(
		ctx context.Context,
		orgID uuid.UUID,
	) (bool, error)
	Update(
		ctx context.Context,
		orgID uuid.UUID,
		params UpdateOrgParams,
	) error

	Delete(ctx context.Context, ID uuid.UUID) error
}

type orgMemberRepo interface {
	Create(
		ctx context.Context,
		member CreateMemberParams,
	) error

	GetByID(
		ctx context.Context,
		memberID uuid.UUID,
	) (models.OrgMember, error)
	GetForAccountAndOrg(
		ctx context.Context,
		accountID, organizationID uuid.UUID,
	) (models.OrgMember, error)

	ExistsByID(
		ctx context.Context,
		memberID uuid.UUID,
	) (bool, error)
	ExistsForAccountAndOrg(
		ctx context.Context,
		accountID, organizationID uuid.UUID,
	) (bool, error)

	Update(
		ctx context.Context,
		memberID uuid.UUID,
		params UpdateMemberParams,
	) error

	Delete(
		ctx context.Context,
		memberID uuid.UUID,
	) error
}

type tombstoneRepo interface {
	OrgMemberIsBuried(ctx context.Context, memberID uuid.UUID) (bool, error)
	BuryOrgMember(ctx context.Context, memberID uuid.UUID) error

	OrganizationIsBuried(ctx context.Context, orgID uuid.UUID) (bool, error)
	BuryOrganization(ctx context.Context, orgID uuid.UUID) error
}

type transaction interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type OrgModule struct {
	org       orgRepo
	member    orgMemberRepo
	tombstone tombstoneRepo
	tx        transaction
}

type OrgModuleDeps struct {
	Org       orgRepo
	Member    orgMemberRepo
	Tombstone tombstoneRepo
	Tx        transaction
}

func NewOrgModule(deps OrgModuleDeps) *OrgModule {
	return &OrgModule{
		org:       deps.Org,
		member:    deps.Member,
		tombstone: deps.Tombstone,
		tx:        deps.Tx,
	}
}

type CreateOrgParams struct {
	ID        uuid.UUID `json:"id"`
	Status    string    `json:"status"`
	Name      string    `json:"name"`
	IconKey   *string   `json:"icon_key,omitempty"`
	BannerKey *string   `json:"banner_key,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}

func (m *OrgModule) Create(
	ctx context.Context,
	org CreateOrgParams,
) error {
	exists, err := m.org.Exists(ctx, org.ID)
	if err != nil {
		return err
	}
	if exists {
		return errx.ErrorOrganizationAlreadyExists.Raise(
			fmt.Errorf("organization with id %s already exists", org.ID),
		)
	}

	bury, err := m.tombstone.OrganizationIsBuried(ctx, org.ID)
	if err != nil {
		return err
	}
	if bury {
		return errx.ErrorOrganizationDeleted.Raise(
			fmt.Errorf("organization with id %s is already deleted", org.ID),
		)
	}

	return m.org.Create(ctx, org)
}

func (m *OrgModule) Get(
	ctx context.Context,
	orgID uuid.UUID,
) (models.Organization, error) {
	return m.org.Get(ctx, orgID)
}

func (m *OrgModule) GetByIDs(
	ctx context.Context,
	ids []uuid.UUID,
) ([]models.Organization, error) {
	return m.org.GetListByIDs(ctx, ids)
}

type UpdateOrgParams struct {
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	IconKey   *string   `json:"icon_key"`
	BannerKey *string   `json:"banner_key"`
	Version   int32     `json:"version"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (m *OrgModule) Update(
	ctx context.Context,
	orgID uuid.UUID,
	params UpdateOrgParams,
) error {
	org, err := m.org.Get(ctx, orgID)
	if errors.Is(err, errx.ErrorOrganizationNotExists) {
		buried, err := m.tombstone.OrganizationIsBuried(ctx, orgID)
		if err != nil {
			return err
		}
		if buried {
			return errx.ErrorOrganizationDeleted.Raise(
				fmt.Errorf("organization with id %s is already deleted", orgID),
			)
		}
	}
	if err != nil {
		return err
	}

	if params.Version <= org.Version {
		return nil
	}

	return m.org.Update(ctx, orgID, params)
}

func (m *OrgModule) Delete(
	ctx context.Context,
	organizationID uuid.UUID,
) error {
	buried, err := m.tombstone.OrganizationIsBuried(ctx, organizationID)
	if err != nil {
		return err
	}
	if buried {
		return errx.ErrorOrganizationDeleted.Raise(
			fmt.Errorf("organization with id %s is already deleted", organizationID),
		)
	}

	return m.tx.Transaction(ctx, func(ctx context.Context) error {
		if err := m.tombstone.BuryOrganization(ctx, organizationID); err != nil {
			return err
		}

		return m.org.Delete(ctx, organizationID)
	})
}
