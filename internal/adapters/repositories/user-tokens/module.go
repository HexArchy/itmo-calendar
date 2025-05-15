package usertokens

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hexarchy/itmo-calendar/internal/config"
)

// Secrets holds secret values for user tokens.
type Secrets struct {
	JWTSecret string `yaml:"jwt_secret"`
}

var Module = fx.Module("repo-user-tokens",
	fx.Provide(func(loader *config.Loader) (*Secrets, error) {
		var s Secrets
		return &s, loader.Unmarshal("secret", &s)
	}),
	fx.Provide(func(db *pgxpool.Pool, secrets *Secrets, logger *zap.Logger) (*Repository, error) {
		return New(db, secrets.JWTSecret, logger)
	}),
)
