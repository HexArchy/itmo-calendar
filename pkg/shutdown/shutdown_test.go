package shutdown

import (
	"context"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBase(t *testing.T) {
	Reset()
	defer Reset()

	i := 0
	cb := &Callback{
		Name: "test base callback",
		Fn: func() error {
			i++
			return nil
		},
	}

	AddCallback(cb)
	AddCallback(cb)
	Shutdown()

	assert.NotEqual(t, 2, i)

	err := Wait(&Config{
		Delay:           0,
		WaitTimeout:     time.Second,
		CallbackTimeout: time.Second,
	})

	assert.Nil(t, err)
	assert.Equal(t, 2, i)
}

func TestBaseWithContext(t *testing.T) {
	Reset()
	defer Reset()

	ctx := WithContext(context.Background())
	canceled := false

	done := make(chan struct{})
	go func() {
		<-ctx.Done()
		canceled = true
		close(done)
	}()

	Shutdown()
	err := Wait(&Config{
		Delay:           0,
		WaitTimeout:     time.Second,
		CallbackTimeout: time.Second,
	})

	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Error("Context was not canceled in time")
	}

	assert.Nil(t, err)
	assert.True(t, canceled)
}

func TestForceStop(t *testing.T) {
	Reset()
	defer Reset()

	cb := &Callback{
		Name: "test slow callback",
		Fn: func() error {
			time.Sleep(5 * time.Second)
			return nil
		},
	}
	AddCallback(cb)

	Shutdown()

	go func() {
		time.Sleep(100 * time.Millisecond)
		globalHandler.forceCh <- syscall.SIGINT
	}()

	err := Wait(&Config{
		Delay:           0,
		WaitTimeout:     time.Second * 5,
		CallbackTimeout: time.Second * 5,
	})

	assert.Equal(t, ErrForceShutdown, err)
}

func TestStopTimeout(t *testing.T) {
	Reset()
	defer Reset()

	cb := &Callback{
		Name: "test slow callback",
		Fn: func() error {
			time.Sleep(5 * time.Second)
			return nil
		},
	}
	AddCallback(cb)
	Shutdown()

	err := Wait(&Config{
		Delay:           0,
		WaitTimeout:     time.Millisecond,
		CallbackTimeout: time.Millisecond,
	})

	assert.Equal(t, ErrTimeoutExceeded, err)
}

func TestCallbackTimeout(t *testing.T) {
	Reset()
	defer Reset()

	cb := &Callback{
		Name: "test timeout callback",
		FnCtx: func(ctx context.Context) error {
			select {
			case <-time.After(2 * time.Second):
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		},
	}
	AddCallback(cb)
	Shutdown()

	err := Wait(&Config{
		Delay:           0,
		WaitTimeout:     5 * time.Second,
		CallbackTimeout: 100 * time.Millisecond,
	})

	assert.NotNil(t, err)
	assert.Equal(t, ErrTimeoutExceeded, err)
}

func TestIsShuttingDown(t *testing.T) {
	Reset()
	defer Reset()

	assert.False(t, IsShuttingDown(), "Should not be shutting down initially")

	Shutdown()
	time.Sleep(10 * time.Millisecond)

	assert.True(t, IsShuttingDown(), "Should be shutting down after Shutdown call")
}

func TestAdd(t *testing.T) {
	Reset()
	defer Reset()

	executed := false
	Add("test-ctx-callback", func(ctx context.Context) error {
		executed = true
		return nil
	})

	Shutdown()

	err := Wait(&Config{
		Delay:           0,
		WaitTimeout:     time.Second,
		CallbackTimeout: time.Second,
	})

	assert.Nil(t, err)
	assert.True(t, executed)
}
