package logger

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

type LogEnv string

const (
	LogDev     LogEnv = "development"
	LogStaging LogEnv = "staging"
	LogProd    LogEnv = "production"
)

func InitLogger(env LogEnv) zerolog.Logger {
	var logLevel zerolog.Level
	var writer io.Writer = os.Stderr

	switch env {
	case LogDev:
		logLevel = zerolog.DebugLevel
		writer = zerolog.ConsoleWriter{Out: os.Stderr}
	case LogStaging:
		logLevel = zerolog.InfoLevel
	case LogProd:
		logLevel = zerolog.WarnLevel
	}

	zerolog.SetGlobalLevel(logLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	return zerolog.New(writer).With().Timestamp().Logger()
}
