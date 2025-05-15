package sendschedule

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	sendscheduleuc "github.com/hexarchy/itmo-calendar/internal/use-cases/send-schedule"
	"github.com/hexarchy/itmo-calendar/pkg/rabbitmq"
)

var Module = fx.Module(
	"worker-send-schedule",
	fx.Provide(
		func(rabbit *rabbitmq.Client, uc *sendscheduleuc.UseCase, cfg *rabbitmq.Config, logger *zap.Logger) *Worker {
			return New(rabbit, uc, cfg.Queues.SendScheduleQueue, logger)
		},
	),
)
