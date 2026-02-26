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
	Create(
		ctx context.Context,
		params models.Organization,
	) error
	GetOrganization(
		ctx context.Context,
		orgID uuid.UUID,
	) (models.Organization, error)
	Update(
		ctx context.Context,
		orgID uuid.UUID,
		params UpdateParams,
	) error
	UpdateOrgStatus(
		ctx context.Context,
		orgID uuid.UUID,
		status string,
		version int32,
		updatedAt time.Time,
	) error
	Delete(ctx context.Context, ID uuid.UUID) error

	CreateOrgMember(
		ctx context.Context,
		member models.OrgMember,
	) error
	UpdateOrgMember(ctx context.Context, memberID uuid.UUID, params UpdateMemberParams) error
	DeleteOrgMember(ctx context.Context, memberID uuid.UUID) error

	UpdatePlaceStatusForOrg(
		ctx context.Context,
		organizationID uuid.UUID,
		status string,
	) error

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
