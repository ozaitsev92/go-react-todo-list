// Package logger builds a structured JSON logger on top of log/slog.
package logger

import (
	"io"
	"log/slog"
	"strings"
)

// New returns a JSON logger writing to w at level (see ParseLevel).
func New(level string, w io.Writer) *slog.Logger {
	return slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{Level: ParseLevel(level)}))
}

// ParseLevel maps debug, info, warn or error to a slog.Level (default info).
func ParseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// Discard returns a logger that drops everything.
func Discard() *slog.Logger {
	return slog.New(slog.DiscardHandler)
}
