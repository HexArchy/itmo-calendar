package http

import (
	"go.uber.org/fx"

	"github.com/hexarchy/itmo-calendar/internal/config"
)

var Module = fx.Module("http",
	fx.Provide(func(loader *config.Loader) (*Config, error) {
		var cfg Config
		return &cfg, loader.Unmarshal("http_server", &cfg)
	}),
	fx.Provide(NewServer),
)
