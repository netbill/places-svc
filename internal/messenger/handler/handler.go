package handler

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/models"
	"github.com/netbill/places-svc/internal/core/modules/organization"
)

type Handler struct {
	modules Modules
}

type Modules struct {
	organization organizationSvc
}

func New(modules Modules) *Handler {
	return &Handler{
		modules: modules,
	}
}

type organizationSvc interface {
	Create(
		ctx context.Context,
		organization models.Organization,
	) error
	Update(
		ctx context.Context,
		organizationID uuid.UUID,
		params organization.UpdateParams,
	) error
	Delete(
		ctx context.Context,
		organizationID uuid.UUID,
	) error

	CreateOrgMember(ctx context.Context, member models.OrgMember) error
	UpdateOrgMember(ctx context.Context, memberID uuid.UUID, params organization.UpdateMemberParams) error
	DeleteOrgMember(ctx context.Context, ID uuid.UUID) error
}
