package caldav

import "go.uber.org/fx"

var Module = fx.Module("repo-caldav",
	fx.Provide(New),
)
