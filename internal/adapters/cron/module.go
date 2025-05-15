package cron

import (
	"go.uber.org/fx"

	"github.com/hexarchy/itmo-calendar/pkg/rabbitmq"
)

var Module = fx.Module("adapter-cron",
	fx.Provide(func(client *rabbitmq.Client, cfg *rabbitmq.Config) *Adapter {
		return New(client, cfg.Queues.CronProcessScheduleQueue, cfg.Queues.SendScheduleQueue)
	}),
)
