package cli

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/alecthomas/kingpin"
	"github.com/netbill/logium"
	"github.com/netbill/places-svc/cmd"
	"github.com/netbill/places-svc/cmd/migrations"
	"github.com/sirupsen/logrus"
)

func Run(args []string) bool {
	cfg := cmd.LoadConfig()

	logium.SetLevel(logrus.DebugLevel)

	log := logium.New()

	lvl, err := logrus.ParseLevel(cfg.Log.Level)
	if err != nil {
		lvl = logrus.InfoLevel
		log.WithField("bad_level", cfg.Log.Level).Warn("unknown log level, fallback to info")
	}
	log.SetLevel(lvl)

	switch {
	case cfg.Log.Format == "json":
		log.SetFormatter(&logrus.JSONFormatter{})
	default:
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}

	var (
		service = kingpin.New("chains-auth", "")
		runCmd  = service.Command("run", "run command")

		serviceCmd     = runCmd.Command("service", "run service")
		migrateCmd     = service.Command("migrate", "migrate command")
		migrateUpCmd   = migrateCmd.Command("up", "migrate storage up")
		migrateDownCmd = migrateCmd.Command("down", "migrate storage down")
	)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup

	c, err := service.Parse(args[1:])
	if err != nil {
		log.WithError(err).Error("failed to parse arguments")
		return false
	}

	switch c {
	case serviceCmd.FullCommand():
		cmd.StartServices(ctx, cfg, log, &wg)
	case migrateUpCmd.FullCommand():
		err = migrations.MigrateUp(cfg.Database.SQL.URL)
	case migrateDownCmd.FullCommand():
		err = migrations.MigrateDown(cfg.Database.SQL.URL)
	default:
		log.Errorf("unknown command %s", c)
		return false
	}
	if err != nil {
		log.WithError(err).Error("failed to exec cmd")
		return false
	}

	wgch := make(chan struct{})
	go func() {
		wg.Wait()
		close(wgch)
	}()

	select {
	case <-ctx.Done():
		log.Printf("Interrupt signal received: %v", ctx.Err())
		<-wgch
	case <-wgch:
		log.Print("All services stopped")
	}

	return true
}
