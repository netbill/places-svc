package handler

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core"
	"github.com/netbill/places-svc/pkg/log"
)

type Handler struct {
	log     *log.Logger
	modules Modules
}

type Modules struct {
	Org organizationSvc
}

func New(log *log.Logger, modules Modules) *Handler {
	return &Handler{
		log:     log,
		modules: modules,
	}
}

type organizationSvc interface {
	Create(
		ctx context.Context,
		params core.CreateOrgParams,
	) error
	Update(
		ctx context.Context,
		organizationID uuid.UUID,
		params core.UpdateOrgParams,
	) error
	Delete(
		ctx context.Context,
		organizationID uuid.UUID,
	) error

	CreateMember(ctx context.Context, member core.CreateMemberParams) error
	UpdateMember(ctx context.Context, memberID uuid.UUID, params core.UpdateMemberParams) error
	DeleteMember(ctx context.Context, ID uuid.UUID) error
}
