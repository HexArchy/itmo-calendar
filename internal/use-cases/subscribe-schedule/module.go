package subscribeschedule

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	caldavsvc "github.com/hexarchy/itmo-calendar/internal/services/caldav"
	icalsvc "github.com/hexarchy/itmo-calendar/internal/services/ical"
	schedulessvc "github.com/hexarchy/itmo-calendar/internal/services/schedules"
	userssvc "github.com/hexarchy/itmo-calendar/internal/services/users"
	"github.com/hexarchy/itmo-calendar/pkg/transactions"
)

var Module = fx.Module("uc-subscribe-schedule",
	fx.Provide(func(
		tx *transactions.Runner,
		schedules *schedulessvc.Service,
		users *userssvc.Service,
		ical *icalsvc.Service,
		caldav *caldavsvc.Service,
		logger *zap.Logger,
	) *UseCase {
		return New(tx, schedules, users, ical, caldav, logger)
	}),
)
