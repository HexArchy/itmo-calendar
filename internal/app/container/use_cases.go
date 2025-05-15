package container

import (
	getical "github.com/hexarchy/itmo-calendar/internal/use-cases/get-ical"
	preparesendschedule "github.com/hexarchy/itmo-calendar/internal/use-cases/prepare-send-schedule"
	sendschedule "github.com/hexarchy/itmo-calendar/internal/use-cases/send-schedule"
	subscribeschedule "github.com/hexarchy/itmo-calendar/internal/use-cases/subscribe-schedule"
)

type UseCases struct {
	PrepareSendSchedule *preparesendschedule.UseCase
	SendSchedule        *sendschedule.UseCase
	SubscirbeSchedule   *subscribeschedule.UseCase
	GetICal             *getical.UseCase
}

func (c *Container) initUseCases() error {
	c.UseCases.PrepareSendSchedule = preparesendschedule.New(
		c.Services.Cron,
		c.Services.Users,
		c.Logger,
	)

	c.UseCases.SendSchedule = sendschedule.New(
		c.Services.Schedules,
		c.Services.Users,
		c.Services.ICal,
		c.Services.CalDav,
		c.Logger,
	)

	c.UseCases.SubscirbeSchedule = subscribeschedule.New(
		c.Services.Schedules,
		c.Services.Users,
		c.Services.ICal,
		c.Services.CalDav,
		c.Logger,
	)

	c.UseCases.GetICal = getical.New(
		c.Services.CalDav,
	)

	return nil
}
