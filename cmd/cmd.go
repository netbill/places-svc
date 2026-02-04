package cmd

import (
	"context"
	"database/sql"
	"sync"

	"github.com/netbill/evebox/box/inbox"
	"github.com/netbill/evebox/box/outbox"
	"github.com/netbill/logium"
	"github.com/netbill/places-svc/internal/core/modules/organization"
	"github.com/netbill/places-svc/internal/core/modules/pclass"
	"github.com/netbill/places-svc/internal/core/modules/place"
	"github.com/netbill/places-svc/internal/core/modules/profile"
	"github.com/netbill/places-svc/internal/messenger"
	"github.com/netbill/places-svc/internal/messenger/inbound"
	"github.com/netbill/places-svc/internal/messenger/outbound"
	"github.com/netbill/places-svc/internal/repository"
	"github.com/netbill/places-svc/internal/rest"
	"github.com/netbill/places-svc/internal/rest/controller"
)

func StartServices(ctx context.Context, cfg Config, log logium.Logger, wg *sync.WaitGroup) {
	run := func(f func()) {
		wg.Add(1)
		go func() {
			f()
			wg.Done()
		}()
	}

	pg, err := sql.Open("postgres", cfg.Database.SQL.URL)
	if err != nil {
		log.Fatal("failed to connect to database", "error", err)
	}

	database := repository.New(pg)

	outBox := outbox.New(pg)
	inBox := inbox.New(pg)

	kafkaOutbound := outbound.New(log, outBox)

	orgSvc := organization.New(database)
	profileSvc := profile.New(database)
	pclasessSvc := pclass.New(database, kafkaOutbound)
	placesSvc := place.New(database, kafkaOutbound)

	kafkaInbound := inbound.New(log, profileSvc, orgSvc)

	ctrl := controller.New(placesSvc, pclasessSvc, log)
	mdll := mdlv.New(cfg.JWT.User.AccessToken.SecretKey, rest.AccountDataCtxKey, log)
	router := rest.New(log, mdll, ctrl)

	kafkaConsumer := messenger.NewConsumer(log, inBox, kafkaInbound, cfg.Kafka.Brokers...)

	kafkaProducer := messenger.NewProducer(log, outBox, cfg.Kafka.Brokers...)

	log.Infof("starting kafka brokers %s", cfg.Kafka.Brokers)

	run(func() { router.Run(ctx, cfg) })

	run(func() { kafkaConsumer.Run(ctx) })

	run(func() { kafkaProducer.Run(ctx) })
}
