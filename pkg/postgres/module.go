package postgres

import (
	"go.uber.org/fx"

	"github.com/hexarchy/itmo-calendar/internal/config"
)

var Module = fx.Module("postgres",
	fx.Provide(func(loader *config.Loader) (*Config, error) {
		var cfg Config
		return &cfg, loader.Unmarshal("postgres", &cfg)
	}),
	fx.Provide(NewPool),
)
