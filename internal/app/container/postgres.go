package container

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hexarchy/itmo-calendar/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

// initPostgresEngine initializes a pgxpool.Pool instance based on the provided configuration.
func (c *Container) initPostgresEngine(ctx context.Context) (*pgxpool.Pool, error) {
	poolCfg, err := createPoolConfig(c.Config)
	if err != nil {
		return nil, errors.Wrap(err, "create pool config")
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, errors.Wrap(err, "new pgxpool")
	}

	return pool, nil
}

// createPoolConfig creates a pgxpool.Config based on the provided Config.
func createPoolConfig(cfg *config.Config) (*pgxpool.Config, error) {
	connStr := buildConnectionURI(cfg.Postgres.Connection)
	poolCfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, errors.Wrap(err, "parse connection string")
	}

	poolCfg.ConnConfig.TLSConfig, err = cfg.Postgres.Connection.TLS.BuildTLSConfig(cfg.HTTPServer.Host)
	if err != nil {
		return nil, errors.Wrap(err, "build TLS config")
	}

	poolCfg.MaxConns = cfg.Postgres.Pool.MaxConnections
	poolCfg.MinConns = cfg.Postgres.Pool.MinConnections
	poolCfg.MaxConnLifetime = cfg.Postgres.Pool.MaxConnectionLifetime
	poolCfg.MaxConnIdleTime = cfg.Postgres.Pool.MaxConnectionIdleTime
	poolCfg.HealthCheckPeriod = cfg.Postgres.Pool.HealthCheckPeriod
	poolCfg.ConnConfig.ConnectTimeout = cfg.Postgres.ConnectTimeout

	if cfg.Postgres.StatementTimeout > 0 {
		if poolCfg.ConnConfig.RuntimeParams == nil {
			poolCfg.ConnConfig.RuntimeParams = make(map[string]string)
		}
		poolCfg.ConnConfig.RuntimeParams["statement_timeout"] = cfg.Postgres.StatementTimeout.String()
	}

	return poolCfg, nil
}

// buildConnectionURI constructs a URI string based on connection configuration.
func buildConnectionURI(conn config.PostgresConnection) string {
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
