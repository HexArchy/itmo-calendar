package shutdown

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var globalHandler *handler

// init initializes the global shutdown handler and sets up signal handling.
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
			fmt.Printf("Signal received: %v\n", sig)

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

	log.Println("Shutdown signal received, initiating graceful shutdown...")

	done := make(chan struct{})
	timer := time.NewTimer(cfg.WaitTimeout)
	defer timer.Stop()

	callbacks := globalHandler.getAllCallbacks()

	shutdownStartTime := time.Now()
	var execErr error
	go func() {
		if cfg.Delay > 0 {
			log.Printf("Waiting %v before shutting down...", cfg.Delay)
			time.Sleep(cfg.Delay)
		}

		// Execute callbacks in reverse order (last added, first executed).
		for i := len(callbacks) - 1; i >= 0; i-- {
			cb := callbacks[i]

			log.Printf("Executing shutdown callback: %s", cb.Name)

			if err := executeCallback(cb, cfg.CallbackTimeout); err != nil {
				log.Printf("Shutdown callback '%s' failed: %v", cb.Name, err)
				if execErr == nil {
					execErr = err
				}
			} else {
				log.Printf("Shutdown callback '%s' completed successfully", cb.Name)
			}
		}

		close(done)
	}()

	secondSignalCh := make(chan os.Signal, 1)

	// Only start listening for force shutdown after a grace period
	// to avoid interpreting the same signal as both normal and force shutdown.
	go func() {
		// Wait for a short time to avoid capturing the same signal twice.
		time.Sleep(250 * time.Millisecond)

		// Now start listening for force signals.
		signal.Notify(secondSignalCh, syscall.SIGINT, syscall.SIGTERM)

		// Forward any signals to the force channel.
		for sig := range secondSignalCh {
			// Ensure it's actually a second signal, not just latency from the first.
			if time.Since(shutdownStartTime) > 500*time.Millisecond {
				globalHandler.forceCh <- sig
				return
			}
		}
	}()

	// Ensure we stop listening for signals when done.
	defer signal.Stop(secondSignalCh)

	select {
	case <-done:
		log.Println("Graceful shutdown completed")
		return execErr
	case <-globalHandler.forceCh:
		log.Println("Force shutdown signal received")
		return ErrForceShutdown
	case <-timer.C:
		log.Println("Shutdown timeout exceeded")
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
