package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	nethttp "net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/embersyndicate/support/internal/mongo"
	"github.com/go-redis/redis/v8"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	mongod "go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/natefinch/lumberjack.v2"
)

// basics initializes the following
// loadConfig - parses environment variables and applies them to a struct
// loadLogger - takes in a configuration and intializes a logrus logger
// loadDB - takes in a configuration and establishes a connection with our datastore, in this application that is mongoDB
// loadRedis - takes in a configuration and establises a connection with our cache, in this application that is Redis
// loadNewrelic - takes in a configuration and configures a NR App to report metrics to NewRelic for monitoring
// loadClient - create a client from the net/http library that is used on all outgoing http requests
func basics(command string) *app {

	app := app{}

	app.cfg, err = loadConfig()
	if err != nil {
		log.Fatalf("failed to load configuration: %s", err)
	}

	app.logger, err = loadLogger(app.cfg, command)
	if err != nil {
		log.Fatalf("failed to load logger: %s", err)
	}

	app.newrelic, err = configNRApplication(app.cfg, app.logger)
	if err != nil {
		app.logger.WithError(err).Fatal("failed to configure NR App")
	}

	app.db, err = makeMongoDB(app.cfg)
	if err != nil {
		app.logger.WithError(err).Fatal("failed to make mongo db connection")
	}

	app.redis = makeRedis(app.cfg)
	if err != nil {
		app.logger.WithError(err).Fatal("failed to configure redis client")
	}

	app.client = &nethttp.Client{
		Timeout:   time.Second * 5,
		Transport: newrelic.NewRoundTripper(nil),
	}
	return &app

}

func loadLogger(cfg config, command string) (*logrus.Logger, error) {
	logger = logrus.New()

	logger.SetOutput(ioutil.Discard)

	logger.AddHook(&writerHook{
		Writer:    os.Stdout,
		LogLevels: logrus.AllLevels,
	})

	logger.AddHook(&writerHook{
		Writer: &lumberjack.Logger{
			Filename: fmt.Sprintf("logs/error-%s.log", command),
			MaxSize:  10,
			Compress: false,
		},
		LogLevels: []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
		},
	})

	logger.AddHook(&writerHook{
		Writer: &lumberjack.Logger{
			Filename:   fmt.Sprintf("logs/info-%s.log", command),
			MaxBackups: 3,
			MaxSize:    10,
			Compress:   false,
		},
		LogLevels: []logrus.Level{
			logrus.InfoLevel,
		},
	})

	level, err := logrus.ParseLevel(cfg.Log.Level)
	if err != nil {
		return logger, errors.Wrap(err, "failed to configure log level")
	}

	logger.SetLevel(level)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// if cfg.Env == "production" {
	// 	// logrus.SetFormatter(nrlogrusplugin.ContextFormatter{})
	// }

	return logger, nil
}

type writerHook struct {
	Writer    io.Writer
	LogLevels []logrus.Level
}

func (w *writerHook) Fire(entry *logrus.Entry) error {

	line, err := entry.String()
	if err != nil {
		return err
	}
	_, err = w.Writer.Write([]byte(line))
	return err
}

func (w *writerHook) Levels() []logrus.Level {
	return w.LogLevels
}

func configNRApplication(cfg config, logger *logrus.Logger) (app *newrelic.Application, err error) {

	opts := []newrelic.ConfigOption{}
	opts = append(opts, newrelic.ConfigFromEnvironment())
	// opts = append(opts, newrelic.ConfigInfoLogger(logger.Writer()))

	app, err = newrelic.NewApplication(opts...)
	if err != nil {
		return nil, err
	}

	err = app.WaitForConnection(time.Second * 20)

	return

}

func makeMongoDB(cfg config) (*mongod.Database, error) {

	q := url.Values{}
	q.Set("maxIdleTimeMS", strconv.FormatInt(int64(time.Second*10), 10))
	q.Set("connectTimeoutMS", strconv.FormatInt(int64(time.Second*4), 10))
	q.Set("serverSelectionTimeoutMS", strconv.FormatInt(int64(time.Second*4), 10))
	q.Set("socketTimeoutMS", strconv.FormatInt(int64(time.Second*4), 10))
	c := &url.URL{
		Scheme:   "mongodb",
		Host:     fmt.Sprintf("%s:%d", cfg.Mongo.Host, cfg.Mongo.Port),
		User:     url.UserPassword(cfg.Mongo.User, cfg.Mongo.Pass),
		Path:     fmt.Sprintf("/%s", cfg.Mongo.Name),
		RawQuery: q.Encode(),
	}

	if cfg.Env == development {
		fmt.Println(c.String())
	}

	mc, err := mongo.Connect(context.TODO(), c)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongo, sleep and continue")
	}

	err = mc.Ping(context.TODO(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping mongo, sleep and continue")
	}

	return mc.Database(cfg.Mongo.Name), nil

}

func makeRedis(cfg config) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:               fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		MaxRetries:         5,
		IdleTimeout:        time.Second * 10,
		IdleCheckFrequency: time.Second * 5,
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		logger.WithError(err).Fatal("failed to ping redis server")
	}

	return redisClient
}
