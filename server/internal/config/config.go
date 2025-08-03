package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type (
	ConfEnv        string
	Configurations struct {
		Port                string
		Env                 ConfEnv
		ShutdownGracePeriod time.Duration
		DatabaseURL         string
	}
)

const (
	ConfEnvDev     ConfEnv = "development"
	ConfEnvStaging ConfEnv = "staging"
	ConfEnvProd    ConfEnv = "production"
)

var AppConfigurations Configurations

func init() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	log.Info().Msgf("Current working directory: %s", dir)

	if err := godotenv.Load(".env"); err != nil {
		log.Warn().Err(err).Msg("could not load .env file")
		panic(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	shutdownGracePeriod := os.Getenv("SHUTDOWN_GRACE_PERIOD")
	if shutdownGracePeriod == "" {
		shutdownGracePeriod = "10"
	}

	databaseURL := os.Getenv("DATABASE_URL")

	var confEnv ConfEnv
	switch env {
	case "development":
		confEnv = ConfEnvDev
	case "staging":
		confEnv = ConfEnvStaging
	case "production":
		confEnv = ConfEnvProd
	default:
		confEnv = ConfEnvDev
	}

	parsedShutdownGP, err := strconv.Atoi(shutdownGracePeriod)
	if err != nil {
		log.Warn().Err(err).Msg("Invalid SHUTDOWN_GRACE_PERIOD, using default 10 seconds")
		parsedShutdownGP = 10
	}

	AppConfigurations = Configurations{
		Port:                port,
		Env:                 confEnv,
		ShutdownGracePeriod: time.Second * time.Duration(parsedShutdownGP),
		DatabaseURL:         databaseURL,
	}

	log.Info().Msg("config is loaded")
}
