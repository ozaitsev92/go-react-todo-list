// Command app runs the todo API server.
package main

import (
	"log/slog"
	"os"

	"github.com/ozaitsev92/tododdd/config"
	"github.com/ozaitsev92/tododdd/internal/app"
	"github.com/ozaitsev92/tododdd/pkg/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("startup failed", "err", err)
		os.Exit(1)
	}

	log := logger.New(cfg.LogLevel, os.Stdout)
	slog.SetDefault(log)

	if err := app.Run(cfg, log); err != nil {
		log.Error("app stopped with error", "err", err)
		os.Exit(1)
	}
}
