package app

import (
	"context"
	"golab/internal/platform/config"
	"golab/internal/platform/database"
	platformlogger "golab/internal/platform/logger"
	platformmail "golab/internal/platform/mail"
	"golab/internal/platform/queue"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Infra struct {
	DB         *pgxpool.Pool
	Rabbit     *queue.RabbitMQ
	SMTPClient *platformmail.SMTPClient
	Logger     *slog.Logger
}

func bootstrapInfra(ctx context.Context, cfg config.Config) (*Infra, error) {
	logger := platformlogger.New(cfg.Environment)

	db, err := database.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("database bootstrap failed", slog.Any("error", err))
		return nil, err
	}

	rabbit, err := queue.NewRabbitMQ(cfg.RabbitMQURL, cfg.RabbitMQExchange)
	if err != nil {
		logger.Error("rabbitmq bootstrap failed", slog.Any("error", err))
		db.Close()
		return nil, err
	}

	if err := rabbit.DeclareQueueAndBind(cfg.RabbitMQActivationQueue, cfg.RabbitMQActivationRoutingKey); err != nil {
		logger.Error("rabbitmq queue declare and bind failed",
			slog.String("queue", cfg.RabbitMQActivationQueue),
			slog.String("routing_key", cfg.RabbitMQActivationRoutingKey),
			slog.Any("error", err),
		)
		_ = rabbit.Close()
		db.Close()
		return nil, err
	}

	return &Infra{
		DB:     db,
		Rabbit: rabbit,
		SMTPClient: platformmail.NewSMTPClient(platformmail.SMTPConfig{
			Host:     cfg.SMTPHost,
			Port:     cfg.SMTPPort,
			Username: cfg.SMTPUsername,
			Password: cfg.SMTPPassword,
		}),
		Logger: logger,
	}, nil
}
