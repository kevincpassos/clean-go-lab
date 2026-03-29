package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"golab/internal/modules/user/usecase"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Handler struct {
	logger  *slog.Logger
	useCase *usecase.UserUseCase
}

func NewHandler(useCase *usecase.UserUseCase, logger *slog.Logger) *Handler {
	if logger == nil {
		logger = slog.Default()
	}

	return &Handler{useCase: useCase, logger: logger}
}

func (h *Handler) StartActivationEmailWorker(ctx context.Context, deliveries <-chan amqp.Delivery) {
	for {
		select {
		case <-ctx.Done():
			return
		case delivery, ok := <-deliveries:
			log := h.deliveryLogger(delivery)

			if !ok {
				h.logger.Warn("worker delivery channel closed")
				return
			}

			if err := h.handleActivationEmail(ctx, delivery.Body); err != nil {
				log.Error("worker activation email failed", slog.Any("error", err))
				requeue := !delivery.Redelivered
				if err := delivery.Nack(false, requeue); err != nil {
					log.Error("worker nack failed", slog.Any("error", err))
				}
				continue
			}

			if err := delivery.Ack(false); err != nil {
				log.Error("worker ack failed", slog.Any("error", err))
			}
		}
	}
}

func (h *Handler) handleActivationEmail(ctx context.Context, payload []byte) error {
	var msg usecase.AccountActivationEmailMessage
	if err := json.Unmarshal(payload, &msg); err != nil {
		return fmt.Errorf("worker: unmarshal activation payload: %w", err)
	}

	if err := h.useCase.SendActivationEmail(ctx, usecase.SendActivationEmailInput{
		Email: msg.Email,
		Name:  msg.Name,
	}); err != nil {
		return fmt.Errorf("worker: send activation email: %w", err)
	}

	return nil
}

func (h *Handler) deliveryLogger(delivery amqp.Delivery) *slog.Logger {
	return h.logger.With(
		slog.String("exchange", delivery.Exchange),
		slog.String("routing_key", delivery.RoutingKey),
		slog.String("message_id", delivery.MessageId),
		slog.String("content_type", delivery.ContentType),
		slog.Bool("redelivered", delivery.Redelivered),
	)
}
