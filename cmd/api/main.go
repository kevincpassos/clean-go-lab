package main

import (
	"context"
	"golab/internal/app"
	"golab/internal/platform/config"
	platformlogger "golab/internal/platform/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()
	bootstrapLogger := platformlogger.New(cfg.Environment)

	container, err := app.NewContainerWithConfig(ctx, cfg)
	if err != nil {
		bootstrapLogger.Error("api container bootstrap failed", "error", err)
		os.Exit(1)
	}
	defer container.Close()

	logger := container.Infra.Logger

	if err := app.StartWorkers(ctx, container); err != nil {
		logger.Error("api worker startup failed", "error", err)
		os.Exit(1)
	}

	server := app.NewAPIServer(container)

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Error("http server shutdown failed", "error", err)
			return
		}
	}()

	logger.Info("api server running", "addr", container.Config.HTTPAddr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("http server failed", "error", err)
		os.Exit(1)
	}
}
