package joblocker

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"

	"github.com/hexarchy/itmo-calendar/pkg/transactions"
)

type Repository struct {
	db *pgxpool.Pool
}

// New returns a new job locker repository.
func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// Lock tries to acquire a lock for the given jobName.
// Returns true if lock acquired, false if already locked.
func (r *Repository) Lock(ctx context.Context, jobName string) (bool, error) {
	const query = `
INSERT INTO job_locks (job_name, locked_at)
VALUES ($1, NOW())
ON CONFLICT (job_name)
DO UPDATE SET locked_at = NOW()
WHERE job_locks.locked_at < NOW() - INTERVAL '1 minute'
    RETURNING job_name
`
	db := transactions.FromContext(ctx, r.db)
	var name string
	err := db.QueryRow(ctx, query, jobName).Scan(&name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, errors.Wrap(err, "acquire lock")
	}

	return true, nil
}

// Unlock releases the lock for the given jobName.
func (r *Repository) Unlock(ctx context.Context, jobName string) error {
	db := transactions.FromContext(ctx, r.db)
	const query = `DELETE FROM job_locks WHERE job_name = $1`
	_, err := db.Exec(ctx, query, jobName)
	return err
}
