package postgres

import (
	"context"
	"golab/internal/modules/user/domain"
	"golab/internal/modules/user/usecase/ports"
	"golab/internal/platform/database"
)

type UserRepository struct {
	db database.DBTX
}

func NewUserRepository(db database.DBTX) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	const query = `
		INSERT INTO users (name, email)
		VALUES ($1, $2)
		RETURNING id, name, email, created_at, updated_at
	`

	var created domain.User
	err := r.db.QueryRow(ctx, query, user.Name, user.Email).Scan(
		&created.ID,
		&created.Name,
		&created.Email,
		&created.CreatedAt,
		&created.UpdatedAt,
	)
	if err != nil {
		return nil, mapPGError(err)
	}

	return &created, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	const query = `
		SELECT id, name, email, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user domain.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, mapPGError(err)
	}

	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, params ports.PatchUserParams) (*domain.User, error) {
	const query = `
		UPDATE users
		SET
			name = COALESCE($2, name),
			email = COALESCE($3, email)
		WHERE id = $1
		RETURNING id, name, email, created_at, updated_at
	`

	var user domain.User
	err := r.db.QueryRow(ctx, query, params.ID, params.Name, params.Email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, mapPGError(err)
	}

	return &user, nil
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM users WHERE id = $1`

	cmd, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return mapPGError(err)
	}

	if cmd.RowsAffected() == 0 {
		return ports.ErrUserNotFound
	}

	return nil
}
