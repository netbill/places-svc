package organization

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/models"
)

type orgRepo interface {
	Create(
		ctx context.Context,
		params CreateParams,
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
		params UpdateParams,
	) error

	Delete(ctx context.Context, ID uuid.UUID) error
}

type memberRepo interface {
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
