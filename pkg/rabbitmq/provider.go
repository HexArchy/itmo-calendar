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

	maxRetries := cfg.MaxRetries
	if maxRetries == 0 {
		maxRetries = _defaultMaxRetries
	}

	maxReconnectAttempts := cfg.MaxReconnectAttempts
	initialReconnectDelay := cfg.InitialReconnectDelay
	if initialReconnectDelay == 0 {
		initialReconnectDelay = int(_defaultReconnectDelay.Seconds())
	}

	client, err := NewWithConfig(
		context.Background(),
		cfg.BuildDSN(),
		tls,
		logger,
		maxRetries,
		maxReconnectAttempts,
		_defaultReconnectDelay,
	)
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
