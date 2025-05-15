package ical

import "go.uber.org/fx"

var Module = fx.Module("svc-ical",
	fx.Provide(New),
)
