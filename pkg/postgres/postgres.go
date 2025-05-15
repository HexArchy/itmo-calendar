package postgres

import (
	"context"
	"fmt"
	"net/url"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"go.uber.org/fx"
)

// NewPool creates a pgxpool.Pool with fx lifecycle hooks.
func NewPool(lc fx.Lifecycle, cfg *Config) (*pgxpool.Pool, error) {
	connStr := BuildConnectionURI(cfg.Connection)
	poolCfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, errors.Wrap(err, "parse connection string")
	}

	poolCfg.ConnConfig.TLSConfig, err = cfg.Connection.TLS.BuildTLSConfig(cfg.Connection.Hosts)
	if err != nil {
		return nil, errors.Wrap(err, "build TLS config")
	}

	poolCfg.MaxConns = cfg.Pool.MaxConnections
	poolCfg.MinConns = cfg.Pool.MinConnections
	poolCfg.MaxConnLifetime = cfg.Pool.MaxConnectionLifetime
	poolCfg.MaxConnIdleTime = cfg.Pool.MaxConnectionIdleTime
	poolCfg.HealthCheckPeriod = cfg.Pool.HealthCheckPeriod
	poolCfg.ConnConfig.ConnectTimeout = cfg.ConnectTimeout

	if cfg.StatementTimeout > 0 {
		if poolCfg.ConnConfig.RuntimeParams == nil {
			poolCfg.ConnConfig.RuntimeParams = make(map[string]string)
		}
		poolCfg.ConnConfig.RuntimeParams["statement_timeout"] = cfg.StatementTimeout.String()
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		return nil, errors.Wrap(err, "new pgxpool")
	}

	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			pool.Close()
			return nil
		},
	})

	return pool, nil
}

// BuildConnectionURI constructs a PostgreSQL URI from connection config.
func BuildConnectionURI(conn Connection) string {
	encodedPassword := url.QueryEscape(conn.Password)
	return fmt.Sprintf(
		"postgresql://%s:%s@%s/%s?%s",
		conn.Username,
		encodedPassword,
		conn.Hosts,
		conn.Database,
		conn.Additional,
	)
}
