package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/heltirj/otus_homeworks/hw12_13_14_15_calendar/internal/app"
	"github.com/heltirj/otus_homeworks/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/heltirj/otus_homeworks/hw12_13_14_15_calendar/internal/server/http"
	"github.com/heltirj/otus_homeworks/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/heltirj/otus_homeworks/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/heltirj/otus_homeworks/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/configs/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	config, err := NewConfig(configFile)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logg := logger.New(config.LogLevel)

	var strg app.Storage

	switch config.StorageType {
	case storage.InMemory:
		strg = memorystorage.New()
	case storage.SQLStorage:
		strg, err = sqlstorage.New(ctx, config.Database.DSN)
		if err != nil {
			log.Fatalf("failed to connect to database: %v", err)
		}

		strg.(*sqlstorage.Storage).Close()
	default:
		log.Fatalf("unknown storage type: %v", config.StorageType)
	}

	calendar := app.New(logg, strg)

	server := internalhttp.NewServer(logg, calendar, config.Service.Host, config.Service.Port)

	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
