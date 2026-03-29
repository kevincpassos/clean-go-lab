package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type TxStarter interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

type TxFunc func(tx pgx.Tx) error

func WithTransaction(ctx context.Context, starter TxStarter, fn TxFunc) (err error) {
	if starter == nil {
		return errors.New("transaction starter is nil")
	}

	if fn == nil {
		return errors.New("transaction function is nil")
	}

	tx, err := starter.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}

		if err != nil {
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil && !errors.Is(rollbackErr, pgx.ErrTxClosed) {
				err = errors.Join(err, fmt.Errorf("rollback transaction: %w", rollbackErr))
			}
		}
	}()

	if err = fn(tx); err != nil {
		return fmt.Errorf("transaction fn: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
