package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/embersyndicate/support/internal/server"
	"github.com/urfave/cli/v2"
)

func serverCommand(c *cli.Context) error {

	basics := basics("server")

	s := server.New(basics.cfg.Server.Port, basics.logger, basics.redis, basics.newrelic)

	serverErrors := make(chan error, 1)

	go func() {
		serverErrors <- s.Run()
	}()

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		basics.logger.WithError(err).Fatal("server encountered an unexpected error and had to quit")
	case sig := <-osSignals:
		basics.logger.WithField("sig", sig).Info("interrupt signal received, starting server shutdown")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		err = s.GracefullyShutdown(ctx)
		if err != nil {
			basics.logger.WithError(err).Fatal("failed to shutdown server")
		}

		basics.logger.Info("server gracefully shutdown successfully")
	}

	return nil

}
