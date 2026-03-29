package postgres

import (
	"errors"
	"golab/internal/modules/user/usecase/ports"
	platformdb "golab/internal/platform/database"

	"github.com/jackc/pgx/v5"
)

func mapPGError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return ports.ErrUserNotFound
	}

	if platformdb.IsUniqueViolation(err) {
		return ports.ErrEmailConflict
	}

	return err
}
