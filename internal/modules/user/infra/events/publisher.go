package events

import (
	"context"
	"fmt"
	"golab/internal/platform/queue"
)

type Publisher struct {
	rabbit     *queue.RabbitMQ
	routingKey string
}

// NewPublisher cria o publisher responsável por enviar eventos
// de usuário para uma routing key específica no RabbitMQ.
func NewPublisher(rabbit *queue.RabbitMQ, routingKey string) *Publisher {
	return &Publisher{
		rabbit:     rabbit,
		routingKey: routingKey,
	}
}

// PublishAccountActivationEmail publica o payload do evento
// de ativação de conta no broker de mensagens.
func (p *Publisher) PublishAccountActivationEmail(ctx context.Context, payload []byte) error {
	if err := p.rabbit.Publish(ctx, p.routingKey, payload); err != nil {
		return fmt.Errorf("publish account activation email: %w", err)
	}
	return nil
}
