package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Environment string

	HTTPAddr string

	DatabaseURL string

	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	SMTPFrom     string

	RabbitMQURL                  string
	RabbitMQExchange             string
	RabbitMQActivationRoutingKey string
	RabbitMQActivationQueue      string
}

var ErrInvalidConfig = errors.New("config: invalid configuration")

func Load() Config {
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbHost := getEnv("DB_HOST", "127.0.0.1")
	dbPort := getEnv("DB_PORT", "5432")
	dbName := getEnv("DB_NAME", "appdb")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")

	databaseURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
		dbSSLMode,
	)

	return Config{
		Environment: getEnv("APP_ENV", "development"),

		HTTPAddr: getEnv("HTTP_ADDR", ":8080"),

		DatabaseURL: databaseURL,

		SMTPHost:     getEnv("SMTP_HOST", "smtp.hostinger.com"),
		SMTPPort:     getEnv("SMTP_PORT", "465"),
		SMTPUsername: getEnv("SMTP_USERNAME", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:     getEnv("SMTP_FROM", getEnv("SMTP_USERNAME", "")),

		RabbitMQURL:                  getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		RabbitMQExchange:             getEnv("RABBITMQ_EXCHANGE", "app.events"),
		RabbitMQActivationRoutingKey: getEnv("RABBITMQ_ACTIVATION_ROUTING_KEY", "user.activation_email"),
		RabbitMQActivationQueue:      getEnv("RABBITMQ_ACTIVATION_QUEUE", "user.activation_email.queue"),
	}
}

func (c Config) Validate() error {
	missing := make([]string, 0, 5)

	if strings.TrimSpace(c.DatabaseURL) == "" {
		missing = append(missing, "DatabaseURL")
	}
	if strings.TrimSpace(c.SMTPUsername) == "" {
		missing = append(missing, "SMTPUsername")
	}
	if strings.TrimSpace(c.SMTPPassword) == "" {
		missing = append(missing, "SMTPPassword")
	}
	if strings.TrimSpace(c.SMTPFrom) == "" {
		missing = append(missing, "SMTPFrom")
	}
	if strings.TrimSpace(c.RabbitMQURL) == "" {
		missing = append(missing, "RabbitMQURL")
	}

	if len(missing) > 0 {
		return fmt.Errorf("%w: missing required fields: %s", ErrInvalidConfig, strings.Join(missing, ", "))
	}

	return nil
}

func getEnv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}
