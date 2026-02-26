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
	UpdateOrgMember(
		ctx context.Context,
		memberID uuid.UUID,
		params UpdateMemberParams,
	) error
	DeleteOrgMember(ctx context.Context, memberID uuid.UUID) error

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
