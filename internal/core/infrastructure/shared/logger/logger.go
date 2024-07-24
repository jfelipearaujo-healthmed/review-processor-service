package logger

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/jfelipearaujo-healthmed/review-processor-service/internal/core/infrastructure/config"
)

const (
	LOG_EMERG   = slog.Level(0)
	LOG_ERR     = slog.Level(3)
	LOG_WARNING = slog.Level(4)
	LOG_NOTICE  = slog.Level(5)
	LOG_DEBUG   = slog.Level(7)
)

func SetupLog(config *config.Config) {
	var level slog.Level
	var handler slog.Handler

	logLevel := "info"

	if config.ApiConfig.IsDevelopment() {
		logLevel = "debug"
	}

	if err := level.UnmarshalText([]byte(logLevel)); err != nil {
		panic(fmt.Errorf("unable to load log level: %v", err))
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	handler = slog.NewJSONHandler(os.Stdout, opts)

	if config.ApiConfig.IsDevelopment() {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	log := slog.New(handler)
	slog.SetDefault(log)
}
