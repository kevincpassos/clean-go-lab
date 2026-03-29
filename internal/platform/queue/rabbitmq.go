package queue

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel

	Exchange string
}

// NewRabbitMQ cria a conexão e o canal com o RabbitMQ,
// além de garantir que o exchange principal exista.
func NewRabbitMQ(url string, exchange string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("rabbitmq dial: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("rabbitmq channel: %w", err)
	}

	r := &RabbitMQ{
		Conn:     conn,
		Channel:  ch,
		Exchange: exchange,
	}

	if err := r.Channel.ExchangeDeclare(
		r.Exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, fmt.Errorf("rabbitmq declare exchange: %w", err)
	}

	return r, nil
}

// DeclareQueueAndBind declara uma fila durável e faz o bind dela
// com a routing key informada no exchange configurado.
func (r *RabbitMQ) DeclareQueueAndBind(queueName, routingKey string) error {
	_, err := r.Channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("rabbitmq declare queue: %w", err)
	}

	if err := r.Channel.QueueBind(
		queueName,
		routingKey,
		r.Exchange,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("rabbitmq bind queue: %w", err)
	}

	return nil
}

// Publish envia uma mensagem para o exchange usando a routing key,
// marcando a entrega como persistente e com content-type JSON.
func (r *RabbitMQ) Publish(ctx context.Context, routingKey string, body []byte) error {
	return r.Channel.PublishWithContext(
		ctx,
		r.Exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
}

// Consume inicia o consumo de mensagens da fila informada.
// O ack manual fica habilitado para o consumidor confirmar o processamento.
func (r *RabbitMQ) Consume(queueName string) (<-chan amqp.Delivery, error) {
	return r.Channel.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
}

// Close encerra canal e conexão com o RabbitMQ de forma segura.
func (r *RabbitMQ) Close() error {
	if r.Channel != nil {
		_ = r.Channel.Close()
	}
	if r.Conn != nil {
		_ = r.Conn.Close()
	}
	return nil
}
