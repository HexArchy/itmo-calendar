package rabbitmq

import (
	"context"
	"time"

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

	reconnectDelay := time.Duration(cfg.InitialReconnectDelay) * time.Second
	if reconnectDelay <= 0 {
		reconnectDelay = _defaultReconnectDelay
	}

	client, err := NewWithConfig(
		context.Background(),
		cfg.BuildDSN(),
		tls,
		logger,
		maxRetries,
		cfg.MaxReconnectAttempts,
		reconnectDelay,
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
