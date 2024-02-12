package main

import (
	"flag"
	"github.com/IlyaZayats/servord/internal/cache"
	"github.com/IlyaZayats/servord/internal/db"
	"github.com/IlyaZayats/servord/internal/handlers"
	"github.com/IlyaZayats/servord/internal/repositories"
	"github.com/IlyaZayats/servord/internal/services"
	"github.com/IlyaZayats/servord/internal/streaming"
	"github.com/sirupsen/logrus"
)

func main() {

	var (
		dbUrl            string
		memcachedUrl     string
		natsStreamingUrl string
		logLevel         string
	)

	flag.StringVar(&dbUrl, "db", "postgres://mylky:mylky@postgres:5432/servord", "database connection url")
	flag.StringVar(&memcachedUrl, "memcached", "memcached:11211", "memcached connection url")
	flag.StringVar(&natsStreamingUrl, "nats-streaming", "nats:4444", "nats-streaming connection url")
	flag.StringVar(&logLevel, "log-level", "info", "log level: panic, fatal, error, warning, info, debug, trace")

	flag.Parse()

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.Panicf("unable to get log level: %v", err)
	}
	logrus.SetLevel(level)

	dbc, err := db.NewPostgresPool(dbUrl)
	if err != nil {
		logrus.Panicf("unable get postgres pool: %v", err)
	}

	mc, err := cache.NewMemcachedCache(memcachedUrl)
	if err != nil {
		logrus.Panicf("unable get memcached: %v", err)
	}

	ns, err := streaming.NewNatsStreamingConnection(natsStreamingUrl)
	if err != nil {
		logrus.Panicf("unable get nats-streaming: %v", err)
	}

	repo, err := repositories.NewPostgresOrderRepository(dbc)
	if err != nil {
		logrus.Panicf("unable init order repo: %v", err)
	}

	service, err := services.NewOrderService(repo)
	if err != nil {
		logrus.Panicf("unable init order service: %v", err)
	}

	handler, err := handlers.NewNatsHandler(service, mc, ns)
	if err != nil {
		logrus.Panicf("unable init order handler: %v", err)
	}

	err = handler.Run()
	if err != nil {
		logrus.Panicf("order handler error: %v", err)
	}

}
