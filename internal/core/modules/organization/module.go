package organization

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/models"
)

type Module struct {
	repo repo
}

func New(repo repo) *Module {
	return &Module{
		repo: repo,
	}
}

type repo interface {
	CreateOrganization(
		ctx context.Context,
		params CreateParams,
	) error
	GetOrganization(
		ctx context.Context,
		orgID uuid.UUID,
	) (models.Organization, error)
	GetOrgsByIDs(
		ctx context.Context,
		ids []uuid.UUID,
	) ([]models.Organization, error)
	ExistsOrganization(
		ctx context.Context,
		orgID uuid.UUID,
	) (bool, error)
	UpdateOrganization(
		ctx context.Context,
		orgID uuid.UUID,
		params UpdateParams,
	) error
	DeleteOrganization(ctx context.Context, ID uuid.UUID) error

	CreateOrgMember(
		ctx context.Context,
		member CreateMemberParams,
	) error
	ExistsOrgMember(
		ctx context.Context,
		memberID uuid.UUID,
	) (bool, error)
	UpdateOrgMember(
		ctx context.Context,
		memberID uuid.UUID,
		params UpdateMemberParams,
	) error
	DeleteOrgMember(ctx context.Context, memberID uuid.UUID) error

	OrgMemberIsBuried(ctx context.Context, memberID uuid.UUID) (bool, error)
	OrganizationIsBuried(ctx context.Context, orgID uuid.UUID) (bool, error)

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
