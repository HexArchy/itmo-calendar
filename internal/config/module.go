package config

import "go.uber.org/fx"

var Module = fx.Module("config",
	fx.Provide(NewLoader),
	fx.Provide(func(l *Loader) (*ITMO, error) {
		var cfg ITMO
		return &cfg, l.Unmarshal("itmo", &cfg)
	}),
	fx.Provide(func(l *Loader) (*Shutdown, error) {
		var cfg Shutdown
		return &cfg, l.Unmarshal("shutdown", &cfg)
	}),
)
