package main

import (
	"context"
	"golab/internal/app"
	"golab/internal/platform/config"
	platformlogger "golab/internal/platform/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()
	bootstrapLogger := platformlogger.New(cfg.Environment)

	container, err := app.NewContainerWithConfig(ctx, cfg)
	if err != nil {
		bootstrapLogger.Error("worker container bootstrap failed", "error", err)
		os.Exit(1)
	}
	defer container.Close()

	logger := container.Infra.Logger

	if err := app.StartWorkers(ctx, container); err != nil {
		logger.Error("worker startup failed", "error", err)
		os.Exit(1)
	}

	<-ctx.Done()
}
