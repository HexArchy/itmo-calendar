package migrations

import (
	"context"
	"embed"

	_ "github.com/jackc/pgx/v5/stdlib" // Import the pgx driver.
	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

//go:embed sql/*.sql
var migrationsFS embed.FS

const _dir = "sql"

// GooseLogger adapts vklog.Logger to goose.Logger interface.
type GooseLogger struct {
	logger *zap.SugaredLogger
}

func (g *GooseLogger) Printf(format string, v ...interface{}) {
	g.logger.Infof(format, v...)
}

func (g *GooseLogger) Fatalf(format string, v ...interface{}) {
	g.logger.Infof(format, v...)
}

func ApplyMigrations(ctx context.Context, logger *zap.Logger, dbString string) error {
	goose.SetBaseFS(migrationsFS)
	goose.SetLogger(&GooseLogger{logger: logger.Sugar()})
	err := goose.SetDialect(string(goose.DialectPostgres))
	if err != nil {
		return errors.Wrap(err, "set dialect")
	}

	db, err := goose.OpenDBWithDriver(string(goose.DialectPostgres), dbString)
	if err != nil {
		return errors.Wrap(err, "open db")
	}
	defer db.Close()

	err = goose.UpContext(ctx, db, _dir)
	if err != nil {
		return errors.Wrap(err, "apply migrations")
	}

	return nil
}
