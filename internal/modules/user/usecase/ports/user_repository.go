package ports

import (
	"context"
	"errors"
	"golab/internal/modules/user/domain"
)

var (
	ErrUserNotFound  = errors.New("users.repository: user not found")
	ErrEmailConflict = errors.New("users.repository: email already exists")
)

type PatchUserParams struct {
	ID    int64
	Name  *string
	Email *string
}

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	GetByID(ctx context.Context, id int64) (*domain.User, error)
	Update(ctx context.Context, params PatchUserParams) (*domain.User, error)
	Delete(ctx context.Context, id int64) error
}
