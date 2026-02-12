package caldav

import (
	"github.com/emersion/go-webdav/caldav"
	"go.uber.org/fx"
	"go.uber.org/zap"

	getical "github.com/hexarchy/itmo-calendar/internal/use-cases/get-ical"
	subscribeschedule "github.com/hexarchy/itmo-calendar/internal/use-cases/subscribe-schedule"
)

// Module provides CalDAV handler components for FX dependency injection.
var Module = fx.Module(
	"caldav-handler",
	fx.Provide(func(getICal *getical.UseCase, logger *zap.Logger) *Backend {
		return NewBackend(getICal, logger)
	}),
	fx.Provide(func(
		getICal *getical.UseCase,
		subscribe *subscribeschedule.UseCase,
		logger *zap.Logger,
	) *AuthMiddleware {
		return NewAuthMiddleware(getICal, subscribe, logger)
	}),
	fx.Provide(func(b *Backend) *caldav.Handler {
		return &caldav.Handler{
			Backend: b,
			Prefix:  "/caldav",
		}
	}),
)
