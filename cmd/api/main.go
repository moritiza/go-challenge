package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/moritiza/go-challenge/internal/config"
	"github.com/moritiza/go-challenge/internal/service"
	redisstorage "github.com/moritiza/go-challenge/internal/storage/redis"
	transporthttp "github.com/moritiza/go-challenge/internal/transport/http"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := run(logger); err != nil {
		logger.Error("service exited with error", "error", err)

		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	store := redisstorage.NewStorage(
		cfg.Redis.Addr,
		cfg.Redis.Password,
		cfg.Redis.DB,
		cfg.Maintenance.CleanupBatchSize,
	)
	defer func() {
		if err := store.Close(); err != nil {
			logger.Error("failed to close redis connection", "error", err)
		}
	}()

	pingCtx, cancel := context.WithTimeout(
		context.Background(),
		cfg.Redis.PingTimeout,
	)
	defer cancel()

	if err := store.Ping(pingCtx); err != nil {
		return fmt.Errorf("connect to redis: %w", err)
	}

	logger.Info("connected to redis", "addr", cfg.Redis.Addr)

	svc := service.NewEstimationService(
		store,
		cfg.Domain.MembershipTTL,
	)

	handler := transporthttp.NewHandler(svc, logger)
	router := transporthttp.NewRouter(handler)

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	return serveWithGracefulShutdown(
		srv,
		svc,
		cfg.Server.ShutdownTimeout,
		cfg.Maintenance.CleanupInterval,
		logger,
	)
}

func serveWithGracefulShutdown(
	srv *http.Server,
	svc *service.EstimationService,
	shutdownTimeout time.Duration,
	cleanupInterval time.Duration,
	logger *slog.Logger,
) error {
	serverErrs := make(chan error, 1)

	go func() {
		logger.Info("starting http server", "addr", srv.Addr)

		if err := srv.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			serverErrs <- fmt.Errorf("http server: %w", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(stop)

	cleanupCtx, cleanupCancel := context.WithCancel(context.Background())
	defer cleanupCancel()

	go runCleanupWorker(
		cleanupCtx,
		svc,
		cleanupInterval,
		logger,
	)

	select {
	case err := <-serverErrs:
		return err

	case sig := <-stop:
		logger.Info("shutdown signal received", "signal", sig.String())
	}

	cleanupCancel()

	ctx, cancel := context.WithTimeout(
		context.Background(),
		shutdownTimeout,
	)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("graceful shutdown failed: %w", err)
	}

	logger.Info("server stopped gracefully")

	return nil
}

func runCleanupWorker(
	ctx context.Context,
	svc *service.EstimationService,
	interval time.Duration,
	logger *slog.Logger,
) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	logger.Info("cleanup worker started", "interval", interval)

	for {
		select {
		case <-ctx.Done():
			logger.Info("cleanup worker stopped")

			return

		case <-ticker.C:
			if err := svc.CleanupAllExpired(ctx); err != nil {
				logger.Error("cleanup failed", "error", err)
			}
		}
	}
}
