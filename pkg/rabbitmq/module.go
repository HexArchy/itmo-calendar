package rabbitmq

import (
	"go.uber.org/fx"

	"github.com/hexarchy/itmo-calendar/internal/config"
)

var Module = fx.Module("rabbitmq",
	fx.Provide(func(loader *config.Loader) (*Config, error) {
		var cfg Config
		return &cfg, loader.Unmarshal("rabbitmq", &cfg)
	}),
	fx.Provide(NewClient),
)
