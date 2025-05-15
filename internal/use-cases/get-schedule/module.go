package getschedule

import (
	"go.uber.org/fx"

	caldavsvc "github.com/hexarchy/itmo-calendar/internal/services/caldav"
	icalsvc "github.com/hexarchy/itmo-calendar/internal/services/ical"
)

var Module = fx.Module("uc-get-schedule",
	fx.Provide(func(caldav *caldavsvc.Service, ical *icalsvc.Service) *UseCase {
		return New(caldav, ical)
	}),
)
