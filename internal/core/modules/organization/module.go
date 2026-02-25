package organization

import (
	"context"
	"time"

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
		params models.Organization,
	) (models.Organization, error)
	UpdateOrganization(
		ctx context.Context,
		orgID uuid.UUID,
		params UpdateParams,
	) (models.Organization, error)
	UpdateOrgStatus(
		ctx context.Context,
		orgID uuid.UUID,
		status string,
		updatedAt time.Time,
	) (models.Organization, error)
	DeleteOrganization(ctx context.Context, ID uuid.UUID) error

	CreateOrgMember(
		ctx context.Context,
		member models.OrgMember,
	) (models.OrgMember, error)
	DeleteOrgMember(ctx context.Context, memberID uuid.UUID) error

	UpdatePlaceStatusForOrg(
		ctx context.Context,
		organizationID uuid.UUID,
		status string,
	) error

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
