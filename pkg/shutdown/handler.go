package shutdown

import (
	"os"
	"sync"
	"sync/atomic"
	"syscall"
)

const (
	statusRunning uint32 = 0
	statusClosing uint32 = 1
)

// handler manages the shutdown process.
type handler struct {
	// signal channel for shutdown signals.
	signalCh chan os.Signal

	// force channel for forced shutdown signals.
	forceCh chan os.Signal

	// mutex for protecting callbacks.
	mu sync.RWMutex

	// user callbacks.
	callbacks []*Callback

	// system callbacks (context cancellation, etc.).
	sysCallbacks []*Callback

	// atomic status flag for shutdown state.
	status uint32
}

// newHandler creates a new shutdown handler.
func newHandler() *handler {
	return &handler{
		signalCh: make(chan os.Signal, 1),
		forceCh:  make(chan os.Signal, 1),
	}
}

// addCallback adds a user callback to be executed during shutdown.
func (h *handler) addCallback(cb *Callback) {
	if cb == nil {
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	h.callbacks = append(h.callbacks, cb)
}

// addSystemCallback adds a system callback to be executed during shutdown.
func (h *handler) addSystemCallback(cb *Callback) {
	if cb == nil {
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	h.sysCallbacks = append(h.sysCallbacks, cb)
}

// getAllCallbacks returns all callbacks (user and system) in execution order.
func (h *handler) getAllCallbacks() []*Callback {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// Create a new slice to avoid race conditions.
	allCallbacks := make([]*Callback, 0, len(h.callbacks)+len(h.sysCallbacks))

	// Copy user callbacks.
	allCallbacks = append(allCallbacks, h.callbacks...)

	// Add system callbacks (will be executed after user callbacks).
	allCallbacks = append(allCallbacks, h.sysCallbacks...)

	return allCallbacks
}

// triggerShutdown sends a signal to initiate the shutdown process.
func (h *handler) triggerShutdown() {
	// Set the shutting down flag first to prevent race conditions.
	h.setShuttingDown()

	// Send signal to indicate shutdown has been triggered.
	select {
	case h.signalCh <- syscall.SIGINT:
		// Signal sent successfully.
	default:
		// Channel is full or closed, no need to send.
	}
}

// setShuttingDown marks the handler as shutting down.
func (h *handler) setShuttingDown() {
	atomic.StoreUint32(&h.status, statusClosing)
}

// isShuttingDown returns true if shutdown has been initiated.
func (h *handler) isShuttingDown() bool {
	return atomic.LoadUint32(&h.status) == statusClosing
}
