package container

import (
	"github.com/hexarchy/itmo-calendar/internal/adapters/cron"
	itmoschedule "github.com/hexarchy/itmo-calendar/internal/adapters/itmo-schedule"
	itmotokens "github.com/hexarchy/itmo-calendar/internal/adapters/itmo-tokens"
	"github.com/hexarchy/itmo-calendar/internal/adapters/repositories/caldav"
	joblocker "github.com/hexarchy/itmo-calendar/internal/adapters/repositories/job-locker"
	usertokens "github.com/hexarchy/itmo-calendar/internal/adapters/repositories/user-tokens"
	"github.com/hexarchy/itmo-calendar/internal/adapters/repositories/users"
)

type Adapters struct {
	ITMOSchedule *itmoschedule.Client
	ITMOTokens   *itmotokens.Client

	Cron *cron.Adapter

	UserTokens *usertokens.Repository
	Users      *users.Repository
	JobLocker  *joblocker.Repository
	CalDav     *caldav.Repository
}

func (c *Container) initAdapters() error {
	c.Adapters.ITMOSchedule = itmoschedule.New(
		c.Config.ITMO.BaseURL,
	)
	c.Adapters.ITMOTokens = itmotokens.New(
		c.Config.ITMO.ClientID,
		c.Config.ITMO.RedirectURI,
		c.Config.ITMO.ProviderURL,
		c.Logger,
	)
	c.Adapters.UserTokens = usertokens.New(
		c.Infra.Postgres,
		c.Config.Secrets.JWTSecret,
		c.Logger,
	)
	c.Adapters.Users = users.New(
		c.Infra.Postgres,
	)
	c.Adapters.JobLocker = joblocker.New(
		c.Infra.Postgres,
	)
	c.Adapters.Cron = cron.New(
		c.Infra.RabbitMQ,
		c.Config.RabbitMQ.Queues.CronProcessScheduleQueue,
		c.Config.RabbitMQ.Queues.SendScheduleQueue,
	)
	c.Adapters.CalDav = caldav.New(
		c.Infra.Postgres,
	)

	return nil
}
