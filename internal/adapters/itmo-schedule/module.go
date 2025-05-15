package itmoschedule

import (
	"go.uber.org/fx"

	"github.com/hexarchy/itmo-calendar/internal/config"
)

var Module = fx.Module("adapter-itmo-schedule",
	fx.Provide(func(cfg *config.ITMO) *Client {
		return New(cfg.BaseURL, cfg.InsecureSkipVerify)
	}),
)
