package messenger

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/eventbox"
	"github.com/netbill/evtypes"
	"github.com/netbill/places-svc/pkg/log"
)

type handlers interface {
	OrganizationCreated(ctx context.Context, event eventbox.InboxEvent) error
	OrganizationUpdated(ctx context.Context, event eventbox.InboxEvent) error
	OrganizationDeleted(ctx context.Context, event eventbox.InboxEvent) error

	OrgMemberCreated(ctx context.Context, event eventbox.InboxEvent) error
	OrgMemberUpdated(ctx context.Context, event eventbox.InboxEvent) error
	OrgMemberDeleted(ctx context.Context, event eventbox.InboxEvent) error
}

func NewInboxWorker(
	logger *log.Logger,
	inbox eventbox.Inbox,
	cfg eventbox.InboxWorkerConfig,
	handlers handlers,
) *eventbox.InboxWorker {
	id := uuid.New().String()

	worker := eventbox.NewInboxWorker(id, logger, inbox, cfg)

	worker.Route(evtypes.OrganizationCreatedEvent, handlers.OrganizationCreated)
	worker.Route(evtypes.OrganizationUpdatedEvent, handlers.OrganizationUpdated)
	worker.Route(evtypes.OrganizationDeletedEvent, handlers.OrganizationDeleted)

	worker.Route(evtypes.OrgMemberCreatedEvent, handlers.OrgMemberCreated)
	worker.Route(evtypes.OrgMemberUpdatedEvent, handlers.OrgMemberUpdated)
	worker.Route(evtypes.OrgMemberDeletedEvent, handlers.OrgMemberDeleted)

	return worker
}
