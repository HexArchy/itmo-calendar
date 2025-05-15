package shutdown

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Callback represents a function to be executed during shutdown.
type Callback struct {
	// Name is a descriptive name for the callback for logging purposes.
	Name string

	// Fn is the function to be executed during shutdown.
	// Deprecated: Use FnCtx instead for better timeout handling.
	Fn func() error

	// FnCtx is the context-aware function to be executed during shutdown.
	FnCtx func(ctx context.Context) error
}

// executeCallback runs a callback with timeout.
func executeCallback(callback *Callback, timeout time.Duration) error {
	if callback == nil {
		return errors.New("nil callback")
	}

	if callback.Name == "" {
		return errors.New("unnamed callback")
	}

	// Create a context with the timeout.
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// If the context-aware function is available, use it.
	if callback.FnCtx != nil {
		err := callback.FnCtx(ctx)
		// Convert context deadline errors to our custom ErrTimeoutExceeded.
		if errors.Is(err, context.DeadlineExceeded) {
			return ErrTimeoutExceeded
		}
		return err
	}

	// Otherwise, use the non-context function.
	if callback.Fn == nil {
		return fmt.Errorf("callback '%s' has no function", callback.Name)
	}

	// Run the non-context function with timeout.
	resultCh := make(chan error, 1)
	go func() {
		resultCh <- callback.Fn()
	}()

	select {
	case err := <-resultCh:
		return err
	case <-ctx.Done():
		return ErrTimeoutExceeded
	}
}
