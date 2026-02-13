package transactions

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

// DBTX is the common interface for pgxpool.Pool and pgx.Tx.
type DBTX interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type ctxKey struct{}

// Runner starts transactions.
type Runner struct {
	pool *pgxpool.Pool
}

// NewRunner creates a new Runner.
func NewRunner(pool *pgxpool.Pool) *Runner {
	return &Runner{pool: pool}
}

// Run executes fn inside a transaction. Commits on success, rollbacks on error.
func (r *Runner) Run(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "begin transaction")
	}

	txCtx := context.WithValue(ctx, ctxKey{}, tx)

	if err = fn(txCtx); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return errors.Wrap(err, "rollback failed: "+rbErr.Error())
		}
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return errors.Wrap(err, "commit transaction")
	}

	return nil
}

// FromContext returns tx from ctx, or fallback (pool).
func FromContext(ctx context.Context, fallback DBTX) DBTX {
	if tx, ok := ctx.Value(ctxKey{}).(DBTX); ok {
		return tx
	}
	return fallback
}
