package logger

import (
	"github.com/rs/zerolog"
)

type LogEnv string

const (
	LogDev     LogEnv = "development"
	LogStaging LogEnv = "staging"
	LogProd    LogEnv = "production"
)

func InitLogger(env LogEnv) {
	var logLevel zerolog.Level
	switch env {
	case LogDev:
		logLevel = zerolog.DebugLevel
	case LogStaging:
		logLevel = zerolog.InfoLevel
	case LogProd:
		logLevel = zerolog.WarnLevel
	}

	zerolog.SetGlobalLevel(logLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}
