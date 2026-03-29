package ports

import "context"

type EventBus interface {
	PublishAccountActivationEmail(ctx context.Context, payload []byte) error
}
