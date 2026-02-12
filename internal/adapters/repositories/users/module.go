package users

import "go.uber.org/fx"

var Module = fx.Module("repo-users",
	fx.Provide(New),
)
