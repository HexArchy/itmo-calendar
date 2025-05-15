package subscribeschedule

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	caldavsvc "github.com/hexarchy/itmo-calendar/internal/services/caldav"
	icalsvc "github.com/hexarchy/itmo-calendar/internal/services/ical"
	schedulessvc "github.com/hexarchy/itmo-calendar/internal/services/schedules"
	userssvc "github.com/hexarchy/itmo-calendar/internal/services/users"
)

var Module = fx.Module("uc-subscribe-schedule",
	fx.Provide(func(
		schedules *schedulessvc.Service,
		users *userssvc.Service,
		ical *icalsvc.Service,
		caldav *caldavsvc.Service,
		logger *zap.Logger,
	) *UseCase {
		return New(schedules, users, ical, caldav, logger)
	}),
)
