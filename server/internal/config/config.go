package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

// load config : godotenv
// config struct

type (
	ConfEnv        string
	Configurations struct {
		Env                 ConfEnv
		ShutdownGracePeriod time.Duration
	}
)

const (
	ConfEnvDev     ConfEnv = "development"
	ConfEnvStaging ConfEnv = "staging"
	ConfEnvProd    ConfEnv = "production"
)

var AppConfigurations Configurations

func init() {
	godotenv.Load()
	env := os.Getenv("ENV")
	var confEnv ConfEnv
	switch env {
	case "development":
		confEnv = ConfEnvDev
	case "staging":
		confEnv = ConfEnvStaging
	case "production":
		confEnv = ConfEnvProd
	}

	log.Info().Any("config:env", confEnv).Msg("")

	shutdownGracePeriod := os.Getenv("SHUTDOWN_GRACE_PERIOD")

	parsedShutdownGP, err := strconv.Atoi(shutdownGracePeriod)
	if err != nil {
		panic(err)
	}

	log.Info().Str("config:shutdownGracePeriod", shutdownGracePeriod).Msg("")

	AppConfigurations = Configurations{
		Env:                 confEnv,
		ShutdownGracePeriod: time.Second * time.Duration(parsedShutdownGP),
	}

	log.Info().Msg("config is loaded")
}
