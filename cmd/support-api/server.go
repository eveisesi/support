package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/embersyndicate/support/internal/category"
	"github.com/embersyndicate/support/internal/key"
	"github.com/embersyndicate/support/internal/server"
	"github.com/embersyndicate/support/internal/ticket"
	"github.com/embersyndicate/support/internal/token"
	"github.com/embersyndicate/support/internal/user"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/urfave/cli/v2"
)

func serverCommand() *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "Initializes the http server that handle http requests to this application",
		Action: func(c *cli.Context) error {

			basics := basics("server")

			repos := initializeRepositories(basics)

			client := &http.Client{
				Timeout: time.Second * 10,
			}
			client.Transport = newrelic.NewRoundTripper(client.Transport)

			categoryServ := category.New(repos.category)
			keyServ := key.New(basics.logger)
			ticketServ := ticket.New(repos.ticket)
			tokenServ := token.New(keyServ)
			userServ := user.New(client, keyServ, tokenServ, repos.user)

			s := server.New(
				basics.cfg.Server.Port,
				basics.logger,
				basics.redis,
				basics.newrelic,
				categoryServ,
				keyServ,
				ticketServ,
				tokenServ,
				userServ,
			)

			serverErrors := make(chan error, 1)

			go func() {
				serverErrors <- s.Run()
			}()

			osSignals := make(chan os.Signal, 1)
			signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)

			select {
			case err := <-serverErrors:
				basics.logger.WithError(err).Error("server encountered an unexpected error and had to quit")
				txn := basics.newrelic.StartTransaction("server error")
				txn.NoticeError(err)
				txn.End()

				basics.logger.Info("shutting down newrelic application")
				basics.newrelic.Shutdown(time.Second * 5)
				basics.logger.Info("newrelic application shutdown successfully")

				os.Exit(1)

			case sig := <-osSignals:
				basics.logger.WithField("sig", sig).Info("interrupt signal received, starting server shutdown")
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
				defer cancel()

				err = s.GracefullyShutdown(ctx)
				if err != nil {
					basics.logger.WithError(err).Fatal("failed to shutdown server")
				}

				basics.logger.Info("server gracefully shutdown successfully")

				basics.logger.Info("shutting down newrelic application")
				basics.newrelic.Shutdown(time.Second * 5)
				basics.logger.Info("newrelic application shutdown successfully")

			}

			return nil

		},
	}
}
