package logging

import (
	"log/slog"
	"os"
)

func BuildLogger(method string) *slog.Logger {
	level := slog.LevelInfo
	if os.Getenv("ENV") != "prod" {
		level = slog.LevelDebug
	} else {
		level = slog.LevelError
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})).With("method", method)
	return logger
}
