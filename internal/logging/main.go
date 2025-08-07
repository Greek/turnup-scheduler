package logging

import (
	"log/slog"
	"os"
)

func BuildLogger(method string) *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})).With("method", method)

	return logger
}
