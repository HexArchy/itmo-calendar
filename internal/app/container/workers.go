package container

import sendschedule "github.com/hexarchy/itmo-calendar/internal/handlers/workers/send-schedule"

type Workers struct {
	RabbitMQ *RabbitMQWorkers
}

type RabbitMQWorkers struct {
	SendSchedule *sendschedule.Worker
}

func (c *Container) initWorkers() error {
	c.Workers.RabbitMQ = &RabbitMQWorkers{
		SendSchedule: sendschedule.New(
			c.Infra.RabbitMQ,
			c.UseCases.SendSchedule,
			c.Config.RabbitMQ.Queues.SendScheduleQueue,
			c.Logger,
		),
	}

	return nil
}
