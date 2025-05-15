package schedules

import (
	"go.uber.org/fx"

	itmoschedule "github.com/hexarchy/itmo-calendar/internal/adapters/itmo-schedule"
	itmotokens "github.com/hexarchy/itmo-calendar/internal/adapters/itmo-tokens"
	usertokens "github.com/hexarchy/itmo-calendar/internal/adapters/repositories/user-tokens"
)

var Module = fx.Module(
	"svc-schedules",
	fx.Provide(
		func(schedule *itmoschedule.Client, tokens *itmotokens.Client, userTokens *usertokens.Repository) *Service {
			return New(schedule, tokens, userTokens)
		},
	),
)
