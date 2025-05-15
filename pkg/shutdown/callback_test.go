package shutdown

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExecuteCallback(t *testing.T) {
	tests := []struct {
		name        string
		callback    *Callback
		timeout     time.Duration
		expectError bool
		errorType   error
	}{
		{
			name:        "nil callback",
			callback:    nil,
			timeout:     time.Second,
			expectError: true,
		},
		{
			name:        "unnamed callback",
			callback:    &Callback{FnCtx: func(ctx context.Context) error { return nil }},
			timeout:     time.Second,
			expectError: true,
		},
		{
			name:        "nil functions",
			callback:    &Callback{Name: "test"},
			timeout:     time.Second,
			expectError: true,
		},
		{
			name: "successful context function",
			callback: &Callback{
				Name:  "test",
				FnCtx: func(ctx context.Context) error { return nil },
			},
			timeout:     time.Second,
			expectError: false,
		},
		{
			name: "successful regular function",
			callback: &Callback{
				Name: "test",
				Fn:   func() error { return nil },
			},
			timeout:     time.Second,
			expectError: false,
		},
		{
			name: "error in context function",
			callback: &Callback{
				Name:  "test",
				FnCtx: func(ctx context.Context) error { return errors.New("test error") },
			},
			timeout:     time.Second,
			expectError: true,
		},
		{
			name: "error in regular function",
			callback: &Callback{
				Name: "test",
				Fn:   func() error { return errors.New("test error") },
			},
			timeout:     time.Second,
			expectError: true,
		},
		{
			name: "timeout in context function",
			callback: &Callback{
				Name: "test",
				FnCtx: func(ctx context.Context) error {
					select {
					case <-time.After(2 * time.Second):
						return nil
					case <-ctx.Done():
						return ctx.Err()
					}
				},
			},
			timeout:     100 * time.Millisecond,
			expectError: true,
			errorType:   ErrTimeoutExceeded,
		},
		{
			name: "timeout in regular function",
			callback: &Callback{
				Name: "test",
				Fn:   func() error { time.Sleep(2 * time.Second); return nil },
			},
			timeout:     100 * time.Millisecond,
			expectError: true,
			errorType:   ErrTimeoutExceeded,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := executeCallback(tt.callback, tt.timeout)

			if tt.expectError {
				assert.Error(t, err, "Expected error but got nil")
			} else {
				assert.NoError(t, err, "Expected no error but got: %v", err)
			}

			if tt.errorType != nil {
				assert.ErrorIs(t, err, tt.errorType, "Expected error type %v but got: %v", tt.errorType, err)
			}
		})
	}
}

func TestExecuteCallbackWithCancellation(t *testing.T) {
	callbackExecuted := false

	callback := &Callback{
		Name: "cancellable-test",
		FnCtx: func(ctx context.Context) error {
			select {
			case <-time.After(500 * time.Millisecond):
				callbackExecuted = true
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		},
	}

	err := executeCallback(callback, 1*time.Second)
	assert.NoError(t, err)
	assert.True(t, callbackExecuted)

	callbackExecuted = false

	err = executeCallback(callback, 50*time.Millisecond)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrTimeoutExceeded)
	assert.False(t, callbackExecuted)
}
