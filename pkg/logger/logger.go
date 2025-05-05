package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

var AppName = "undefined"

func MustSetupLogger(app, stage string, debug bool, level string) zerolog.Logger {
	zerolog.MessageFieldName = "rest"
	zerolog.LevelFieldName = "severity"
	zerolog.TimestampFieldName = "timestamp"
	zerolog.TimeFieldFormat = time.RFC3339Nano
	AppName = app

	var logs zerolog.Logger

	if debug {
		logs = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	} else {
		logs = log.Output(os.Stderr)
	}

	parsedLvl, err := zerolog.ParseLevel(level)

	if err != nil {
		panic(err)
	}

	log.Logger = logs.Level(parsedLvl).With().Str("service", app).Str("stage", stage).Logger()

	return log.Logger
}
