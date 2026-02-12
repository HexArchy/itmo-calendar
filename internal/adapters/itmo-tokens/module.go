package itmotokens

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/hexarchy/itmo-calendar/internal/config"
)

var Module = fx.Module("adapter-itmo-tokens",
	fx.Provide(func(cfg *config.ITMO, logger *zap.Logger) *Client {
		return New(cfg.ClientID, cfg.RedirectURI, cfg.ProviderURL, logger)
	}),
)
