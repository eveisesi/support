package main

import (
	"log"
	"os"

	nethttp "net/http"

	"github.com/go-redis/redis/v8"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	logger *logrus.Logger
	err    error
)

type app struct {
	cfg      config
	newrelic *newrelic.Application
	logger   *logrus.Logger
	db       *mongo.Database
	redis    *redis.Client
	client   *nethttp.Client
}

func main() {
	app := cli.NewApp()
	app.Name = "Ember Support"
	app.UsageText = "ember-support"
	app.Commands = []*cli.Command{
		serverCommand(),
		testCommand(),
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
