package joblocker

import "go.uber.org/fx"

var Module = fx.Module("repo-job-locker",
	fx.Provide(New),
)
