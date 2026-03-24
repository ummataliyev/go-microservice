package logger

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/thealish/go-microservice/internal/config"
)

// New creates a configured zerolog.Logger based on the logging config and environment.
func New(cfg config.LoggingConfig, env string) zerolog.Logger {
	level := parseLevel(cfg.Level)
	zerolog.SetGlobalLevel(level)

	if strings.ToLower(env) == "development" {
		return zerolog.New(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}).With().Timestamp().Caller().Logger()
	}

	return zerolog.New(os.Stdout).With().Timestamp().Logger()
}

// WithRequestID returns a sub-logger with the request_id field attached.
func WithRequestID(logger zerolog.Logger, requestID string) zerolog.Logger {
	return logger.With().Str("request_id", requestID).Logger()
}

func parseLevel(level string) zerolog.Level {
	switch strings.ToLower(level) {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn", "warning":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel
	}
}
