package mailer

import (
	"context"
	"fmt"
	platformmail "golab/internal/platform/mail"
	"log/slog"
)

type Service struct {
	client *platformmail.SMTPClient
	from   string
	logger *slog.Logger
}

func NewService(client *platformmail.SMTPClient, from string, logger *slog.Logger) *Service {
	if logger == nil {
		logger = slog.Default()
	}

	return &Service{
		client: client,
		from:   from,
		logger: logger,
	}
}

func (s *Service) SendActivationEmail(ctx context.Context, to string, name string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if err := s.client.Send(ctx, platformmail.Message{
		From:    s.from,
		To:      to,
		Subject: "Ativacao de conta",
		Body:    buildActivationMessage(name),
	}); err != nil {
		return err
	}

	return nil
}

func buildActivationMessage(name string) string {
	return fmt.Sprintf("Oi %s,\n\nSeu cadastro foi criado com sucesso.", name)
}
