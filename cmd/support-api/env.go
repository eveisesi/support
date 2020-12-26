package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Mongo struct {
		Host string `envconfig:"MONGO_HOST" required:"true"`
		Port int    `envconfig:"MONGO_PORT" required:"true"`
		User string `envconfig:"MONGO_USER" required:"true"`
		Pass string `envconfig:"MONGO_PASS" required:"true"`
		Name string `envconfig:"MONGO_NAME" required:"true"`
	}

	Redis struct {
		Host string `envconfig:"REDIS_HOST" required:"true"`
		Port uint   `envconfig:"REDIS_PORT" required:"true"`
	}

	Env environment `envconfig:"ENV" required:"true"`

	Developer struct {
		Name string `envconfig:"DEVERLOPER_NAME"`
	}

	Log struct {
		Level string `envconfig:"LOG_LEVEL" default:"info"`
	}

	Server struct {
		Port uint `envconfig:"SERVER_PORT" required:"true"`
	}
}

type environment string

const production environment = "production"
const development environment = "development"

func (e environment) String() string {
	return string(e)
}

var validEnvironments = []environment{production, development}

func (c config) validateEnvironment() bool {
	for _, env := range validEnvironments {
		if c.Env == env {
			return true
		}
	}

	return false
}

func loadConfig() (cfg config, err error) {
	err = godotenv.Load("app.env")
	if err != nil {
		fmt.Println(err)
	}

	err = envconfig.Process("", &cfg)
	if err != nil {
		return config{}, err
	}

	if !cfg.validateEnvironment() {
		return config{}, fmt.Errorf("invalid env %s declared", cfg.Env)
	}

	return

}
