package rabbitmq

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// NewClient creates a rabbitmq.Client with fx lifecycle hooks.
func NewClient(lc fx.Lifecycle, cfg *Config, logger *zap.Logger) (*Client, error) {
	tls, err := cfg.TLS.BuildTLSConfig(cfg.Host)
	if err != nil {
		return nil, errors.Wrap(err, "init rabbitmq tls config")
	}

	client, err := New(context.Background(), cfg.BuildDSN(), tls, logger)
	if err != nil {
		return nil, errors.Wrap(err, "init rabbitmq client")
	}

	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			return client.Close()
		},
	})

	return client, nil
}
