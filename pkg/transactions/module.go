package transactions

import "go.uber.org/fx"

var Module = fx.Module("transactions",
	fx.Provide(NewRunner),
)
