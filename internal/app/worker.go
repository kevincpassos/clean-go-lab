package app

import (
	"context"
	"fmt"
	"log/slog"
)

func StartWorkers(ctx context.Context, container *Container) error {
	deliveries, err := container.Infra.Rabbit.Consume(container.Config.RabbitMQActivationQueue)
	if err != nil {
		container.Infra.Logger.Error("worker consumer start failed",
			slog.String("queue", container.Config.RabbitMQActivationQueue),
			slog.Any("error", err),
		)
		return fmt.Errorf("start activation email consumer: %w", err)
	}

	go container.User.WorkerHandler.StartActivationEmailWorker(ctx, deliveries)
	return nil
}
