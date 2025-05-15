package shutdown

import (
	"time"
)

const (
	_defaultDelay           = 5 * time.Second
	_defaultWaitTimeout     = 10 * time.Second
	_defaultCallbackTimeout = 2 * time.Second
)

// Config contains configuration options for the shutdown process.
type Config struct {
	// Delay is the time to wait before executing callbacks.
	Delay time.Duration `default:"5s"`

	// WaitTimeout is the maximum time to wait for the entire shutdown process.
	WaitTimeout time.Duration `default:"10s"`

	// CallbackTimeout is the maximum time allowed for each callback to complete.
	CallbackTimeout time.Duration `default:"2s"`
}

// DefaultConfig provides the default configuration values.
func DefaultConfig() *Config {
	return &Config{
		Delay:           _defaultDelay,
		WaitTimeout:     _defaultWaitTimeout,
		CallbackTimeout: _defaultCallbackTimeout,
	}
}
