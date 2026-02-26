package app

import (
	"context"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/netbill/awsx"
	"github.com/netbill/eventbox"
	"github.com/netbill/places-svc/internal/repository/pg"

	awscfg "github.com/aws/aws-sdk-go-v2/config"
	eventpg "github.com/netbill/eventbox/pg"
	"github.com/netbill/pgdbx"
	"github.com/netbill/places-svc/internal/bucket"
	"github.com/netbill/places-svc/internal/core/modules/organization"
	"github.com/netbill/places-svc/internal/core/modules/pclass"
	"github.com/netbill/places-svc/internal/core/modules/place"
	"github.com/netbill/places-svc/internal/geogueser"
	"github.com/netbill/places-svc/internal/messenger"
	"github.com/netbill/places-svc/internal/messenger/handler"
	"github.com/netbill/places-svc/internal/messenger/publisher"
	"github.com/netbill/places-svc/internal/repository"
	"github.com/netbill/places-svc/internal/rest"
	"github.com/netbill/places-svc/internal/rest/controller"
	"github.com/netbill/places-svc/internal/rest/middlewares"
	"github.com/netbill/places-svc/internal/tokenmanager"
)

func (a *App) Run(ctx context.Context) error {
	var wg = &sync.WaitGroup{}

	run := func(f func()) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			f()
		}()
	}

	pool, err := a.config.PoolDB(ctx)
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	defer pool.Close()

	a.log.Info("starting application")

	db := pgdbx.NewDB(pool)

	repo := &repository.Repository{
		TransactionSql:   pg.NewTransaction(db),
		PlacesSql:        pg.NewPlacesQ(db),
		PlaceClassesSql:  pg.NewPlaceClassesQ(db),
		OrganizationsSql: pg.NewOrganizationsQ(db),
		OrgMembersSql:    pg.NewOrgMembersQ(db),
	}

	cfg, err := awscfg.LoadDefaultConfig(
		context.Background(),
		awscfg.WithRegion(a.config.S3.Aws.Region),
		awscfg.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				a.config.S3.Aws.AccessKeyID,
				a.config.S3.Aws.SecretAccessKey,
				a.config.S3.Aws.SessionToken,
			),
		),
	)
	if err != nil {
		return fmt.Errorf("load aws config: %w", err)
	}

	s3 := bucket.NewStorage(awsx.New(a.config.S3.Aws.BucketName, cfg), bucket.Config{
		LinkTTL:        a.config.S3.Media.Link.TTL,
		PlaceClassIcon: a.config.S3.Media.Resources.PlaceClass.Icon,
		PlaceIcon:      a.config.S3.Media.Resources.Place.Icon,
		PlaceBanner:    a.config.S3.Media.Resources.Place.Banner,
	})

	outbox := eventpg.NewOutbox(db)
	inbox := eventpg.NewInbox(db)

	producer := messenger.NewProducer(a.log, messenger.ProducerConfig{
		Producer: a.config.Kafka.Identity,
		Brokers:  a.config.Kafka.Brokers,
		PlacesV1: messenger.ProduceKafkaConfig{
			RequiredAcks: a.config.Kafka.Produce.Topics.PlacesV1.RequiredAcks,
			Compression:  a.config.Kafka.Produce.Topics.PlacesV1.Compression,
			Balancer:     a.config.Kafka.Produce.Topics.PlacesV1.Balancer,
			BatchSize:    a.config.Kafka.Produce.Topics.PlacesV1.BatchSize,
			BatchTimeout: a.config.Kafka.Produce.Topics.PlacesV1.BatchTimeout,
		},
		PlaceClassesV1: messenger.ProduceKafkaConfig{
			RequiredAcks: a.config.Kafka.Produce.Topics.PlaceClassesV1.RequiredAcks,
			Compression:  a.config.Kafka.Produce.Topics.PlaceClassesV1.Compression,
			Balancer:     a.config.Kafka.Produce.Topics.PlaceClassesV1.Balancer,
			BatchSize:    a.config.Kafka.Produce.Topics.PlaceClassesV1.BatchSize,
			BatchTimeout: a.config.Kafka.Produce.Topics.PlaceClassesV1.BatchTimeout,
		},
	})
	defer producer.Close()

	outbound := publisher.New(a.config.Kafka.Identity, outbox, producer)

	geo, err := geogueser.New()
	if err != nil {
		return fmt.Errorf("create geography: %w", err)
	}

	orgCore := organization.New(repo)
	placeCore := place.New(repo, s3, outbound, geo)
	classCore := pclass.New(repo, s3, outbound)

	tokenManager := tokenmanager.New(tokenmanager.Config{
		Issuer:   a.config.Auth.Tokens.Issuer,
		AccessSK: a.config.Auth.Tokens.AccountAccess.SecretKey,
	})

	ctrl := controller.New(&controller.Modules{
		Place: placeCore,
		Org:   orgCore,
		Class: classCore,
	})
	mdlv := middlewares.New(tokenManager)
	router := rest.New(mdlv, ctrl)

	run(func() {
		router.Run(ctx, a.log, rest.Config{
			Port:              a.config.Rest.Port,
			ReadTimeout:       a.config.Rest.Timeouts.Read,
			ReadHeaderTimeout: a.config.Rest.Timeouts.ReadHeader,
			WriteTimeout:      a.config.Rest.Timeouts.Write,
			IdleTimeout:       a.config.Rest.Timeouts.Idle,
		})
	})

	outboxWorker := messenger.NewOutboxWorker(a.log, outbox, producer, eventbox.OutboxWorkerConfig{
		Routines:       a.config.Kafka.Outbox.Routines,
		Slots:          a.config.Kafka.Outbox.Slots,
		BatchSize:      a.config.Kafka.Outbox.BatchSize,
		Sleep:          a.config.Kafka.Outbox.Sleep,
		MinNextAttempt: a.config.Kafka.Outbox.MinNextAttempt,
		MaxNextAttempt: a.config.Kafka.Outbox.MaxNextAttempt,
		MaxAttempts:    a.config.Kafka.Outbox.MaxAttempts,
	})
	defer outboxWorker.Clean()

	run(func() {
		outboxWorker.Run(ctx)
	})

	inbound := handler.New(handler.Modules{
		Org: orgCore,
	})

	inboxWorker := messenger.NewInboxWorker(a.log, inbox, eventbox.InboxWorkerConfig{
		Routines:       a.config.Kafka.Inbox.Routines,
		Slots:          a.config.Kafka.Inbox.Slots,
		BatchSize:      a.config.Kafka.Inbox.BatchSize,
		Sleep:          a.config.Kafka.Inbox.Sleep,
		MinNextAttempt: a.config.Kafka.Inbox.MinNextAttempt,
		MaxNextAttempt: a.config.Kafka.Inbox.MaxNextAttempt,
		MaxAttempts:    a.config.Kafka.Inbox.MaxAttempts,
	}, inbound)
	defer inboxWorker.Clean()

	run(func() {
		inboxWorker.Run(ctx)
	})

	consumer := messenger.NewConsumer(a.log, inbox, messenger.ConsumerConfig{
		GroupID:    a.config.Kafka.Identity,
		Brokers:    a.config.Kafka.Brokers,
		MinBackoff: a.config.Kafka.Consume.Backoff.Min,
		MaxBackoff: a.config.Kafka.Consume.Backoff.Max,
		OrganizationsV1: messenger.ConsumeKafkaConfig{
			Instances:     a.config.Kafka.Consume.Topics.OrganizationsV1.Instances,
			MinBytes:      a.config.Kafka.Consume.Topics.OrganizationsV1.MinBytes,
			MaxBytes:      a.config.Kafka.Consume.Topics.OrganizationsV1.MaxBytes,
			MaxWait:       a.config.Kafka.Consume.Topics.OrganizationsV1.MaxWait,
			QueueCapacity: a.config.Kafka.Consume.Topics.OrganizationsV1.QueueCapacity,
		},
		OrgMembersV1: messenger.ConsumeKafkaConfig{
			Instances:     a.config.Kafka.Consume.Topics.OrganizationMembersV1.Instances,
			MinBytes:      a.config.Kafka.Consume.Topics.OrganizationMembersV1.MinBytes,
			MaxBytes:      a.config.Kafka.Consume.Topics.OrganizationMembersV1.MaxBytes,
			MaxWait:       a.config.Kafka.Consume.Topics.OrganizationMembersV1.MaxWait,
			QueueCapacity: a.config.Kafka.Consume.Topics.OrganizationMembersV1.QueueCapacity,
		},
	})
	defer consumer.Close()

	run(func() {
		consumer.Run(ctx)
	})

	wg.Wait()
	return nil
}
