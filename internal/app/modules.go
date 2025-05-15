package app

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hexarchy/itmo-calendar/internal/config"
	"github.com/hexarchy/itmo-calendar/migrations"

	"github.com/pkg/errors"
)

func (a *App) startMigrations(ctx context.Context) error {
	a.Logger.Info("Migration – started")

	dbstring := buildConnectionURI(a.Cfg.Postgres.Connection)
	err := migrations.ApplyMigrations(ctx, a.Logger, dbstring)
	if err != nil {
		return errors.Wrap(err, "apply migrations")
	}

	a.Logger.Info("Migration – completed successfully")
	return nil
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
