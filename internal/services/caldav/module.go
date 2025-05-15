package caldav

import (
	"go.uber.org/fx"

	caldavrepo "github.com/hexarchy/itmo-calendar/internal/adapters/repositories/caldav"
)

var Module = fx.Module("svc-caldav",
	fx.Provide(func(r *caldavrepo.Repository) *Service {
		return New(r)
	}),
)
