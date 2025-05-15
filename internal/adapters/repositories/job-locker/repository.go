package joblocker

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
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
	var name string
	err := r.db.QueryRow(ctx, query, jobName).Scan(&name)
	if err != nil {
		// If no rows returned, lock was not acquired.
		return false, errors.Wrap(err, "acquire lock")
	}

	return true, nil
}

// Unlock releases the lock for the given jobName.
func (r *Repository) Unlock(ctx context.Context, jobName string) error {
	const query = `DELETE FROM job_locks WHERE job_name = $1`
	_, err := r.db.Exec(ctx, query, jobName)
	return err
}
