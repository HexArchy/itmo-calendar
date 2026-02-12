package app

import (
	"go.uber.org/fx"

	cronadapter "github.com/hexarchy/itmo-calendar/internal/adapters/cron"
	itmoschedule "github.com/hexarchy/itmo-calendar/internal/adapters/itmo-schedule"
	itmotokens "github.com/hexarchy/itmo-calendar/internal/adapters/itmo-tokens"
	caldavrepo "github.com/hexarchy/itmo-calendar/internal/adapters/repositories/caldav"
	joblockerrepo "github.com/hexarchy/itmo-calendar/internal/adapters/repositories/job-locker"
	usertokensrepo "github.com/hexarchy/itmo-calendar/internal/adapters/repositories/user-tokens"
	usersrepo "github.com/hexarchy/itmo-calendar/internal/adapters/repositories/users"
	"github.com/hexarchy/itmo-calendar/internal/config"
	caldavhandler "github.com/hexarchy/itmo-calendar/internal/handlers/caldav"
	httpserver "github.com/hexarchy/itmo-calendar/internal/handlers/http"
	api "github.com/hexarchy/itmo-calendar/internal/handlers/http/v1"
	sendscheduleworker "github.com/hexarchy/itmo-calendar/internal/handlers/workers/send-schedule"
	caldavsvc "github.com/hexarchy/itmo-calendar/internal/services/caldav"
	cronsvc "github.com/hexarchy/itmo-calendar/internal/services/cron"
	icalsvc "github.com/hexarchy/itmo-calendar/internal/services/ical"
	schedulessvc "github.com/hexarchy/itmo-calendar/internal/services/schedules"
	userssvc "github.com/hexarchy/itmo-calendar/internal/services/users"
	getical "github.com/hexarchy/itmo-calendar/internal/use-cases/get-ical"
	getschedule "github.com/hexarchy/itmo-calendar/internal/use-cases/get-schedule"
	preparesendschedule "github.com/hexarchy/itmo-calendar/internal/use-cases/prepare-send-schedule"
	sendschedule "github.com/hexarchy/itmo-calendar/internal/use-cases/send-schedule"
	subscribeschedule "github.com/hexarchy/itmo-calendar/internal/use-cases/subscribe-schedule"
	"github.com/hexarchy/itmo-calendar/pkg/logger"
	"github.com/hexarchy/itmo-calendar/pkg/postgres"
	"github.com/hexarchy/itmo-calendar/pkg/rabbitmq"
)

// Module is the root fx module that assembles the entire application.
var Module = fx.Options(
	// config
	config.Module,

	// infra
	logger.Module,
	postgres.Module,
	rabbitmq.Module,

	// adapters
	cronadapter.Module,
	itmoschedule.Module,
	itmotokens.Module,
	caldavrepo.Module,
	joblockerrepo.Module,
	usertokensrepo.Module,
	usersrepo.Module,

	// services
	caldavsvc.Module,
	cronsvc.Module,
	icalsvc.Module,
	schedulessvc.Module,
	userssvc.Module,

	// use-cases
	getical.Module,
	getschedule.Module,
	preparesendschedule.Module,
	sendschedule.Module,
	subscribeschedule.Module,

	// handlers
	api.Module,
	caldavhandler.Module,
	httpserver.Module,
	sendscheduleworker.Module,

	// lifecycle
	fx.Invoke(RunMigrations),
	fx.Invoke(Run),
)
