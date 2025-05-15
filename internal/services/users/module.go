package users

import (
	"go.uber.org/fx"

	usersrepo "github.com/hexarchy/itmo-calendar/internal/adapters/repositories/users"
)

var Module = fx.Module("svc-users",
	fx.Provide(func(r *usersrepo.Repository) *Service {
		return New(r)
	}),
)
