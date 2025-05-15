package shutdown

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandlerAddCallback(t *testing.T) {
	h := newHandler()

	h.addCallback(nil)
	assert.Empty(t, h.callbacks, "Handler should not add nil callbacks")

	cb := &Callback{Name: "test", Fn: func() error { return nil }}
	h.addCallback(cb)
	assert.Len(t, h.callbacks, 1, "Handler should add valid callbacks")
	assert.Equal(t, cb, h.callbacks[0], "Added callback should match original")

	sysCb := &Callback{Name: "system-test", Fn: func() error { return nil }}
	h.addSystemCallback(sysCb)
	assert.Len(t, h.sysCallbacks, 1, "Handler should add system callbacks")
	assert.Equal(t, sysCb, h.sysCallbacks[0], "Added system callback should match original")
}

func TestHandlerGetAllCallbacks(t *testing.T) {
	h := newHandler()

	cb1 := &Callback{Name: "test1", Fn: func() error { return nil }}
	cb2 := &Callback{Name: "test2", Fn: func() error { return nil }}
	h.addCallback(cb1)
	h.addCallback(cb2)

	sysCb1 := &Callback{Name: "system1", Fn: func() error { return nil }}
	sysCb2 := &Callback{Name: "system2", Fn: func() error { return nil }}
	h.addSystemCallback(sysCb1)
	h.addSystemCallback(sysCb2)

	all := h.getAllCallbacks()

	assert.Len(t, all, 4, "Should return all callbacks")

	assert.Equal(t, cb1, all[0], "First regular callback should be first")
	assert.Equal(t, cb2, all[1], "Second regular callback should be second")
	assert.Equal(t, sysCb1, all[2], "First system callback should be third")
	assert.Equal(t, sysCb2, all[3], "Second system callback should be fourth")
}

func TestHandlerShuttingDownFlag(t *testing.T) {
	h := newHandler()

	assert.False(t, h.isShuttingDown(), "Should not be shutting down initially")

	h.setShuttingDown()
	assert.True(t, h.isShuttingDown(), "Should be shutting down after setting flag")
}

func TestHandlerTriggerShutdown(t *testing.T) {
	h := newHandler()

	signalReceived := make(chan struct{})
	go func() {
		<-h.signalCh
		close(signalReceived)
	}()

	h.triggerShutdown()

	<-signalReceived
	assert.True(t, h.isShuttingDown(), "Trigger should set shutting down flag")
}

func TestHandlerConcurrentAccess(t *testing.T) {
	h := newHandler()

	done := make(chan struct{})

	go func() {
		for i := 0; i < 100; i++ {
			h.addCallback(&Callback{
				Name:  "concurrent-test",
				FnCtx: func(context.Context) error { return nil },
			})
		}
		close(done)
	}()

	done2 := make(chan struct{})

	go func() {
		for i := 0; i < 100; i++ {
			h.addSystemCallback(&Callback{
				Name:  "concurrent-system-test",
				FnCtx: func(context.Context) error { return nil },
			})
		}
		close(done2)
	}()

	<-done
	<-done2

	assert.Len(t, h.callbacks, 100, "Should have 100 regular callbacks")
	assert.Len(t, h.sysCallbacks, 100, "Should have 100 system callbacks")

	all := h.getAllCallbacks()
	assert.Len(t, all, 200, "Should have 200 total callbacks")
}
