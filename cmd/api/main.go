package main

import (
	"log/slog"
	"os"

	"github.com/moritiza/go-challenge/internal/config"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	logger.Info("application starting",
		"server_port", cfg.Server.Port,
		"redis_addr", cfg.Redis.Addr,
		"membership_ttl", cfg.Domain.MembershipTTL,
	)
}
