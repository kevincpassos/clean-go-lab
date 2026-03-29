package ports

import "context"

type Mailer interface {
	SendActivationEmail(ctx context.Context, to string, name string) error
}
