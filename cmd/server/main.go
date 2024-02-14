package main

import (
	"context"
	"flag"
	"github.com/IlyaZayats/servord/internal/cache"
	"github.com/IlyaZayats/servord/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	var (
		memcachedUrl string
		listen       string
		logLevel     string
	)
	flag.StringVar(&memcachedUrl, "memcached", "memcached:11211", "memcached connection url")
	flag.StringVar(&listen, "listen", ":8080", "server listen interface")
	flag.StringVar(&logLevel, "log-level", "info", "log level: panic, fatal, error, warning, info, debug, trace")

	flag.Parse()

	ctx := context.Background()

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.Panicf("unable to get log level: %v", err)
	}
	logrus.SetLevel(level)

	mc, err := cache.NewMemcachedCache(memcachedUrl)
	if err != nil {
		logrus.Panicf("unable get memcached: %v", err)
	}

	g := gin.New()

	_, err = handlers.NewOrderHandler(g, mc)
	if err != nil {
		logrus.Panicf("unable build order handlers: %v", err)
	}

	doneC := make(chan error)

	go func() { doneC <- g.Run(listen) }()

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGABRT, syscall.SIGHUP, syscall.SIGTERM)

	childCtx, cancel := context.WithCancel(ctx)
	go func() {
		sig := <-signalChan
		logrus.Debugf("exiting with signal: %v", sig)
		cancel()
	}()

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				doneC <- ctx.Err()
			}
		}
	}(childCtx)

	<-doneC
}
