package messenger

import (
	"context"
	"sync"
	"time"

	"github.com/netbill/evebox/box/inbox"
	"github.com/netbill/evebox/consumer"
	"github.com/netbill/places-svc/internal/messenger/contracts"
	"github.com/segmentio/kafka-go"
)

type handlers interface {
	OrgMemberCreated(
		ctx context.Context,
		event inbox.Event,
	) inbox.EventStatus
	OrgMemberDeleted(
		ctx context.Context,
		event inbox.Event,
	) inbox.EventStatus
	OrgMemberAddedRole(
		ctx context.Context,
		event inbox.Event,
	) inbox.EventStatus
	OrgMemberRemovedRole(
		ctx context.Context,
		event inbox.Event,
	) inbox.EventStatus

	OrgRoleCreated(
		ctx context.Context,
		event inbox.Event,
	) inbox.EventStatus
	OrgRoleDeleted(
		ctx context.Context,
		event inbox.Event,
	) inbox.EventStatus
	OrgRolePermissionsUpdated(
		ctx context.Context,
		event inbox.Event,
	) inbox.EventStatus
	OrgRolesRanksUpdated(
		ctx context.Context,
		event inbox.Event,
	) inbox.EventStatus

	OrganizationCreated(
		ctx context.Context,
		event inbox.Event,
	) inbox.EventStatus
	OrganizationDeleted(
		ctx context.Context,
		event inbox.Event,
	) inbox.EventStatus
	OrganizationActivated(
		ctx context.Context,
		event inbox.Event,
	) inbox.EventStatus
	OrganizationDeactivated(
		ctx context.Context,
		event inbox.Event,
	) inbox.EventStatus
}

func (m *Messenger) RunConsumer(ctx context.Context, handlers handlers) {
	wg := &sync.WaitGroup{}
	run := func(f func()) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			f()
		}()
	}

	orgConsumer := consumer.New(
		consumer.NewConsumerParams{
			Log:  m.log,
			DB:   m.db,
			Name: "profiles-svc-profile-consumer",
			Addr: m.addr,
			OnUnknown: func(ctx context.Context, m kafka.Message, eventType string) error {
				return nil
			},
		},
	)

	orgConsumer.Handle(contracts.OrgMemberCreatedEvent, handlers.OrgMemberCreated)
	orgConsumer.Handle(contracts.OrgMemberDeletedEvent, handlers.OrgMemberDeleted)
	orgConsumer.Handle(contracts.OrgMemberRoleAddedEvent, handlers.OrgMemberAddedRole)
	orgConsumer.Handle(contracts.OrgMemberRoleRemovedEvent, handlers.OrgMemberRemovedRole)

	orgConsumer.Handle(contracts.OrgRoleCreatedEvent, handlers.OrgRoleCreated)
	orgConsumer.Handle(contracts.OrgRoleDeletedEvent, handlers.OrgRoleDeleted)
	orgConsumer.Handle(contracts.OrgRolePermissionsUpdatedEvent, handlers.OrgRolePermissionsUpdated)
	orgConsumer.Handle(contracts.OrgRolesRanksUpdatedEvent, handlers.OrgRolesRanksUpdated)

	orgConsumer.Handle(contracts.OrganizationCreatedEvent, handlers.OrganizationCreated)
	orgConsumer.Handle(contracts.OrganizationDeletedEvent, handlers.OrganizationDeleted)
	orgConsumer.Handle(contracts.OrganizationActivatedEvent, handlers.OrganizationActivated)
	orgConsumer.Handle(contracts.OrganizationDeactivatedEvent, handlers.OrganizationDeactivated)

	inboxer1 := consumer.NewInboxer(consumer.NewInboxerParams{
		Log:        m.log,
		Pool:       m.db,
		Name:       "places-svc-inbox-worker-1",
		BatchSize:  10,
		RetryDelay: 1 * time.Minute,
		MinSleep:   100 * time.Millisecond,
		MaxSleep:   1 * time.Second,
	})

	inboxer1.Handle(contracts.OrgMemberCreatedEvent, handlers.OrgMemberCreated)
	inboxer1.Handle(contracts.OrgMemberDeletedEvent, handlers.OrgMemberDeleted)
	inboxer1.Handle(contracts.OrgMemberRoleAddedEvent, handlers.OrgMemberAddedRole)
	inboxer1.Handle(contracts.OrgMemberRoleRemovedEvent, handlers.OrgMemberRemovedRole)

	inboxer1.Handle(contracts.OrgRoleCreatedEvent, handlers.OrgRoleCreated)
	inboxer1.Handle(contracts.OrgRoleDeletedEvent, handlers.OrgRoleDeleted)
	inboxer1.Handle(contracts.OrgRolePermissionsUpdatedEvent, handlers.OrgRolePermissionsUpdated)
	inboxer1.Handle(contracts.OrgRolesRanksUpdatedEvent, handlers.OrgRolesRanksUpdated)

	inboxer1.Handle(contracts.OrganizationCreatedEvent, handlers.OrganizationCreated)
	inboxer1.Handle(contracts.OrganizationDeletedEvent, handlers.OrganizationDeleted)
	inboxer1.Handle(contracts.OrganizationActivatedEvent, handlers.OrganizationActivated)
	inboxer1.Handle(contracts.OrganizationDeactivatedEvent, handlers.OrganizationDeactivated)

	run(func() {
		orgConsumer.Run(ctx, contracts.PlaceSvcGroup, contracts.OrganizationsTopicV1, m.addr...)
	})

	run(func() {
		inboxer1.Run(ctx)
	})

	wg.Wait()
}
