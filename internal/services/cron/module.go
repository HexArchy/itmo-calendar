package cron

import (
	"go.uber.org/fx"

	cronadapter "github.com/hexarchy/itmo-calendar/internal/adapters/cron"
)

var Module = fx.Module("svc-cron",
	fx.Provide(func(a *cronadapter.Adapter) *Service {
		return New(a)
	}),
)
