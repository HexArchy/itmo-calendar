package shutdown

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIntegrationFullWorkflow(t *testing.T) {
	Reset()
	defer Reset()

	ctx := WithContext(context.Background())

	var (
		mu              sync.Mutex
		callbackCounter int
		contextCanceled bool
	)

	for i := 0; i < 3; i++ {
		Add("test-callback", func(ctx context.Context) error {
			mu.Lock()
			callbackCounter++
			mu.Unlock()
			return nil
		})
	}

	go func() {
		<-ctx.Done()
		mu.Lock()
		contextCanceled = true
		mu.Unlock()
	}()

	assert.False(t, IsShuttingDown(), "Should not be shutting down initially")

	Shutdown()

	assert.True(t, IsShuttingDown(), "Should be shutting down after Shutdown call")
	err := Wait(&Config{
		Delay:           0,
		WaitTimeout:     time.Second,
		CallbackTimeout: time.Second,
	})

	assert.NoError(t, err, "Expected no errors during shutdown")

	mu.Lock()
	assert.Equal(t, 3, callbackCounter, "All callbacks should have executed")
	assert.True(t, contextCanceled, "Context should have been canceled")
	mu.Unlock()
}

func TestReuseAfterReset(t *testing.T) {
	Reset()

	executed1 := false
	Add("test-1", func(ctx context.Context) error {
		executed1 = true
		return nil
	})

	Shutdown()
	Wait(nil)

	assert.True(t, executed1, "First callback should have executed")

	Reset()

	executed2 := false
	Add("test-2", func(ctx context.Context) error {
		executed2 = true
		return nil
	})

	assert.False(t, IsShuttingDown(), "Status should be reset")

	Shutdown()
	Wait(nil)

	assert.True(t, executed2, "Second callback should have executed")
}
