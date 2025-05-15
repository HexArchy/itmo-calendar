package shutdown

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// DefaultSignals are the signals that trigger a shutdown by default.
var DefaultSignals = []os.Signal{syscall.SIGINT, syscall.SIGTERM}

// SignalListener listens for OS signals and notifies when they occur.
type SignalListener interface {
	// Listen starts listening for signals.
	Listen(signals []os.Signal) <-chan os.Signal

	// Stop stops listening for signals.
	Stop()
}

// OSSignalListener is the default implementation of SignalListener.
type OSSignalListener struct {
	signalCh chan os.Signal
	stopOnce sync.Once
}

// NewOSSignalListener creates a new signal listener.
func NewOSSignalListener() *OSSignalListener {
	return &OSSignalListener{
		signalCh: make(chan os.Signal, 1),
	}
}

// Listen starts listening for the specified signals.
func (l *OSSignalListener) Listen(signals []os.Signal) <-chan os.Signal {
	signal.Notify(l.signalCh, signals...)
	return l.signalCh
}

// Stop stops listening for signals.
func (l *OSSignalListener) Stop() {
	l.stopOnce.Do(func() {
		signal.Stop(l.signalCh)
		close(l.signalCh)
	})
}
