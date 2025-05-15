package shutdown

import (
	"errors"
)

// Common errors returned by the shutdown package.
var (
	// ErrTimeoutExceeded is returned when the application fails to shutdown within the specified timeout.
	ErrTimeoutExceeded = errors.New("graceful shutdown failed: timeout exceeded")

	// ErrForceShutdown is returned when a second shutdown signal is received during graceful shutdown.
	ErrForceShutdown = errors.New("graceful shutdown failed: force shutdown occurred")
)
