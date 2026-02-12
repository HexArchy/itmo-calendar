package shutdown

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	_signalGracePeriod = 250 * time.Millisecond
	_forceSignalDelay  = 500 * time.Millisecond
)

var globalHandler *handler

// init initializes the global shutdown handler and sets up signal handling.
//
//nolint:gochecknoinits // init is used here to ensure the shutdown handler is ready before any other code runs.
func init() {
	setupHandler()
}

// setupHandler creates and initializes the global shutdown handler.
func setupHandler() {
	globalHandler = newHandler()
	primarySignalCh := make(chan os.Signal, 1)
	signal.Notify(primarySignalCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for sig := range primarySignalCh {
			slog.Info("Signal received", "signal", sig)

			if globalHandler.isShuttingDown() {
				select {
				case globalHandler.forceCh <- sig:
				default:
				}
			} else {
				globalHandler.setShuttingDown()
				select {
				case globalHandler.signalCh <- sig:
				default:
				}
			}
		}
	}()
}

// AddCallback registers a callback for execution before shutdown.
func AddCallback(cb *Callback) {
	globalHandler.addCallback(cb)
}

// Add registers a context-aware callback for execution before shutdown.
func Add(name string, fn func(ctx context.Context) error) {
	AddCallback(&Callback{
		Name:  name,
		FnCtx: fn,
	})
}

// WithContext creates or extends the given context with cancellation on shutdown.
func WithContext(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	newCtx, cancel := context.WithCancel(ctx)
	globalHandler.addSystemCallback(&Callback{
		Name: "context-cancellation",
		FnCtx: func(_ context.Context) error {
			cancel()
			return nil
		},
	})

	return newCtx
}

// Wait waits for application shutdown.
//
// If a second signal is received during shutdown, ErrForceShutdown is returned.
// If the shutdown process exceeds the configured timeout, ErrTimeoutExceeded is returned.
func Wait(config *Config) error {
	cfg := config
	if cfg == nil {
		cfg = DefaultConfig()
	}

	if !globalHandler.isShuttingDown() {
		<-globalHandler.signalCh
		globalHandler.setShuttingDown()
	}

	//nolint:sloglint // pkg shutdown should not depend on app logger.
	slog.Info("Shutdown signal received, initiating graceful shutdown")

	done := make(chan struct{})
	timer := time.NewTimer(cfg.WaitTimeout)
	defer timer.Stop()

	callbacks := globalHandler.getAllCallbacks()

	shutdownStartTime := time.Now()
	var execErr error
	go func() {
		if cfg.Delay > 0 {
			slog.Info("Waiting before shutting down", "delay", cfg.Delay)
			time.Sleep(cfg.Delay)
		}

		// Execute callbacks in reverse order (last added, first executed).
		for i := len(callbacks) - 1; i >= 0; i-- {
			cb := callbacks[i]

			slog.Info("Executing shutdown callback", "name", cb.Name)

			if err := executeCallback(cb, cfg.CallbackTimeout); err != nil {
				slog.Error("Shutdown callback failed", "name", cb.Name, "error", err)
				if execErr == nil {
					execErr = err
				}
			} else {
				slog.Info("Shutdown callback completed successfully", "name", cb.Name)
			}
		}

		close(done)
	}()

	secondSignalCh := make(chan os.Signal, 1)

	// Only start listening for force shutdown after a grace period
	// to avoid interpreting the same signal as both normal and force shutdown.
	go func() {
		// Wait for a short time to avoid capturing the same signal twice.
		time.Sleep(_signalGracePeriod)

		// Now start listening for force signals.
		signal.Notify(secondSignalCh, syscall.SIGINT, syscall.SIGTERM)

		// Forward any signals to the force channel.
		for sig := range secondSignalCh {
			// Ensure it's actually a second signal, not just latency from the first.
			if time.Since(shutdownStartTime) > _forceSignalDelay {
				globalHandler.forceCh <- sig
				return
			}
		}
	}()

	// Ensure we stop listening for signals when done.
	defer signal.Stop(secondSignalCh)

	select {
	case <-done:
		//nolint:sloglint // pkg shutdown should not depend on app logger.
		slog.Info("Graceful shutdown completed")
		return execErr
	case <-globalHandler.forceCh:
		//nolint:sloglint // pkg shutdown should not depend on app logger.
		slog.Info("Force shutdown signal received")
		return ErrForceShutdown
	case <-timer.C:
		//nolint:sloglint // pkg shutdown should not depend on app logger.
		slog.Info("Shutdown timeout exceeded")
		return ErrTimeoutExceeded
	}
}

// IsShuttingDown returns true if the application is shutting down.
func IsShuttingDown() bool {
	return globalHandler.isShuttingDown()
}

// Shutdown initiates the shutdown process.
func Shutdown() {
	// Set the flag first, then trigger the signal.
	globalHandler.setShuttingDown()

	// Send the signal to the channel.
	select {
	case globalHandler.signalCh <- syscall.SIGINT:
		// Signal sent successfully.
	default:
		// Channel is full or closed, no need to send.
	}
}

// Reset reinitializes the shutdown handler (for testing).
func Reset() {
	// Create a completely new handler with clean state.
	setupHandler()
}
