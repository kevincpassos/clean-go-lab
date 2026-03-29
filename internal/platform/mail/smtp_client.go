package mail

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/smtp"
	"strings"
)

var ErrSMTPConfigInvalid = errors.New("platform.mail: invalid smtp config")

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

type Message struct {
	From    string
	To      string
	Subject string
	Body    string
}

type SMTPClient struct {
	cfg SMTPConfig
}

func NewSMTPClient(cfg SMTPConfig) *SMTPClient {
	return &SMTPClient{cfg: cfg}
}

func (c *SMTPClient) Send(ctx context.Context, message Message) error {
	if err := c.validateConfig(); err != nil {
		return err
	}

	if err := validateMessage(message); err != nil {
		return err
	}

	address := net.JoinHostPort(c.cfg.Host, c.cfg.Port)

	dialer := &tls.Dialer{
		NetDialer: &net.Dialer{},
		Config: &tls.Config{
			ServerName: c.cfg.Host,
			MinVersion: tls.VersionTLS12,
		},
	}

	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, c.cfg.Host)
	if err != nil {
		return err
	}
	defer client.Close()

	auth := smtp.PlainAuth("", c.cfg.Username, c.cfg.Password, c.cfg.Host)
	if err := client.Auth(auth); err != nil {
		return err
	}

	if err := client.Mail(message.From); err != nil {
		return err
	}

	if err := client.Rcpt(message.To); err != nil {
		return err
	}

	writer, err := client.Data()
	if err != nil {
		return err
	}

	if _, err := writer.Write([]byte(buildSMTPMessage(message))); err != nil {
		_ = writer.Close()
		return err
	}

	if err := writer.Close(); err != nil {
		return err
	}

	if err := client.Quit(); err != nil {
		return err
	}

	return nil
}

func (c *SMTPClient) validateConfig() error {
	if strings.TrimSpace(c.cfg.Host) == "" ||
		strings.TrimSpace(c.cfg.Port) == "" ||
		strings.TrimSpace(c.cfg.Username) == "" ||
		strings.TrimSpace(c.cfg.Password) == "" {
		return ErrSMTPConfigInvalid
	}

	return nil
}

func validateMessage(message Message) error {
	if strings.TrimSpace(message.From) == "" ||
		strings.TrimSpace(message.To) == "" ||
		strings.TrimSpace(message.Subject) == "" {
		return ErrSMTPConfigInvalid
	}

	return nil
}

func buildSMTPMessage(message Message) string {
	return strings.Join([]string{
		"From: " + message.From,
		"To: " + message.To,
		"Subject: " + message.Subject,
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"",
		message.Body,
	}, "\r\n")
}
