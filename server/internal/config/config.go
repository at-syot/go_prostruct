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

type Configurations struct {
	Env                 string
	ShutdownGracePeriod time.Duration
}

var AppConfigurations Configurations

func init() {
	godotenv.Load()
	env := os.Getenv("ENV")
	log.Info().Str("config:env", env).Msg("")

	shutdownGracePeriod := os.Getenv("SHUTDOWN_GRACE_PERIOD")

	parsedShutdownGP, err := strconv.Atoi(shutdownGracePeriod)
	if err != nil {
		panic(err)
	}

	log.Info().Str("config:shutdownGracePeriod", shutdownGracePeriod).Msg("")

	AppConfigurations = Configurations{
		Env:                 env,
		ShutdownGracePeriod: time.Second * time.Duration(parsedShutdownGP),
	}

	log.Info().Msg("config is loaded")
}
