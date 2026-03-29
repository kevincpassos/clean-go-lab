package logger

import (
	"log/slog"
	"os"
	"strings"
)

func New(environment string) *slog.Logger {
	env := strings.TrimSpace(strings.ToLower(environment))
	if env == "" {
		env = "development"
	}

	level := slog.LevelInfo
	if env == "development" {
		level = slog.LevelDebug
	}

	options := &slog.HandlerOptions{Level: level}

	if env == "production" {
		return slog.New(slog.NewJSONHandler(os.Stdout, options))
	}

	return slog.New(slog.NewTextHandler(os.Stdout, options))
}
