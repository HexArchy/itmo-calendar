package logger

import (
	"go.uber.org/fx"

	"github.com/hexarchy/itmo-calendar/internal/config"
)

var Module = fx.Module("logger",
	fx.Provide(func(loader *config.Loader) (*Config, error) {
		var cfg Config
		return &cfg, loader.Unmarshal("logger", &cfg)
	}),
	fx.Provide(func(loader *config.Loader) (*AppInfo, error) {
		var app AppInfo
		return &app, loader.Unmarshal("app", &app)
	}),
	fx.Provide(New),
)
