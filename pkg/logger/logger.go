package logger

import (
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"log"
	"os"
	"time"
)

func MustSetupLogger(app, stage string, debug bool, level string) zerolog.Logger {
	zerolog.MessageFieldName = "rest"
	zerolog.LevelFieldName = "severity"
	zerolog.TimestampFieldName = "timestamp"
	zerolog.TimeFieldFormat = time.RFC3339Nano

	var logs zerolog.Logger

	if debug {
		logs = zlog.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	} else {
		logs = zlog.Output(os.Stderr)
	}

	parsedLvl, err := zerolog.ParseLevel(level)

	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	zlog.Logger = logs.Level(parsedLvl).With().Str("service", app).Str("stage", stage).Logger()

	return zlog.Logger
}
