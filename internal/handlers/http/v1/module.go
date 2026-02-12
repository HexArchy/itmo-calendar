package api

import (
	"go.uber.org/fx"

	"github.com/hexarchy/itmo-calendar/internal/handlers/http/v1/gen"
)

var Module = fx.Module("api-v1",
	fx.Provide(NewOgenHandler),
	fx.Provide(func(h *OgenHandler) (gen.Handler, error) { return h, nil }),
	fx.Provide(func(h gen.Handler) (*gen.Server, error) {
		return gen.NewServer(h, gen.WithPathPrefix("/api/v1"))
	}),
)
