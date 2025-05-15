package getical

import (
	"go.uber.org/fx"

	caldavsvc "github.com/hexarchy/itmo-calendar/internal/services/caldav"
)

var Module = fx.Module("uc-get-ical",
	fx.Provide(func(caldav *caldavsvc.Service) *UseCase {
		return New(caldav)
	}),
)
