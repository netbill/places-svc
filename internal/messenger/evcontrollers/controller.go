package evcontrollers

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/organization"
	"github.com/netbill/places-svc/pkg/log"
)

type OrgController struct {
	log     *log.Logger
	modules Modules
}

type Modules struct {
	Org orgCore
}

func New(log *log.Logger, modules Modules) *OrgController {
	return &OrgController{
		log:     log,
		modules: modules,
	}
}

type orgCore interface {
	Create(
		ctx context.Context,
		params organization.CreateParams,
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

	CreateMember(
		ctx context.Context,
		member organization.CreateMemberParams,
	) error
	UpdateMember(
		ctx context.Context,
		memberID uuid.UUID,
		params organization.UpdateMemberParams,
	) error
	DeleteMember(
		ctx context.Context,
		ID uuid.UUID,
	) error
}
