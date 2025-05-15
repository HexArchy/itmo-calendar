package api

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DBHealthChecker implements health check using pgxpool.
type DBHealthChecker struct {
	pool *pgxpool.Pool
}

// NewDBHealthChecker creates a new DBHealthChecker.
func NewDBHealthChecker(pool *pgxpool.Pool) *DBHealthChecker {
	return &DBHealthChecker{
		pool: pool,
	}
}

// Ping checks database connectivity.
func (h *DBHealthChecker) Ping(ctx context.Context) error {
	return h.pool.Ping(ctx)
}
