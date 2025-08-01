package config

import (
	"os"
	"path"
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
	}
)

const (
	ConfEnvDev     ConfEnv = "development"
	ConfEnvStaging ConfEnv = "staging"
	ConfEnvProd    ConfEnv = "production"
)

var AppConfigurations Configurations

func init() {
	// Explicitly declare env path,
	// due to we run - air - from root of go workspace
	// REMARK: maybe only when on @development
	wd, _ := os.Getwd()
	envPath := path.Join(wd, "/server/.env")
	godotenv.Load(envPath)
	log.Info().Msgf("load env from %s successful", envPath)

	port := os.Getenv("PORT")
	env := os.Getenv("ENV")
	shutdownGracePeriod := os.Getenv("SHUTDOWN_GRACE_PERIOD")

	var confEnv ConfEnv
	switch env {
	case "development":
		confEnv = ConfEnvDev
	case "staging":
		confEnv = ConfEnvStaging
	case "production":
		confEnv = ConfEnvProd
	}

	parsedShutdownGP, err := strconv.Atoi(shutdownGracePeriod)
	if err != nil {
		panic(err)
	}

	AppConfigurations = Configurations{
		Port:                port,
		Env:                 confEnv,
		ShutdownGracePeriod: time.Second * time.Duration(parsedShutdownGP),
	}

	log.Info().Msg("config is loaded")
}
