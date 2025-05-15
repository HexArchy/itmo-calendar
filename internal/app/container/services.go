package container

import (
	"github.com/hexarchy/itmo-calendar/internal/services/caldav"
	"github.com/hexarchy/itmo-calendar/internal/services/cron"
	"github.com/hexarchy/itmo-calendar/internal/services/ical"
	"github.com/hexarchy/itmo-calendar/internal/services/schedules"
	"github.com/hexarchy/itmo-calendar/internal/services/users"
)

type Services struct {
	Schedules *schedules.Service
	Users     *users.Service
	ICal      *ical.Service
	Cron      *cron.Service
	CalDav    *caldav.Service
}

func (c *Container) initServices() error {
	c.Services.Schedules = schedules.New(
		c.Adapters.ITMOSchedule,
		c.Adapters.ITMOTokens,
		c.Adapters.UserTokens,
	)

	c.Services.Users = users.New(
		c.Adapters.Users,
	)

	c.Services.Cron = cron.New(
		c.Adapters.Cron,
	)

	c.Services.ICal = ical.New()

	c.Services.CalDav = caldav.New(
		c.Adapters.CalDav,
	)

	return nil
}
