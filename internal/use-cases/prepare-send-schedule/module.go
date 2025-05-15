package preparesendschedule

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	cronsvc "github.com/hexarchy/itmo-calendar/internal/services/cron"
	userssvc "github.com/hexarchy/itmo-calendar/internal/services/users"
)

var Module = fx.Module("uc-prepare-send-schedule",
	fx.Provide(func(cron *cronsvc.Service, users *userssvc.Service, logger *zap.Logger) *UseCase {
		return New(cron, users, logger)
	}),
)
