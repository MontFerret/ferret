package debugger

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

type trackingContext struct {
	context.Context

	mu      sync.Mutex
	done    chan struct{}
	nextID  int
	stopped int
	pending map[int]func()
}

func (c *trackingContext) Done() <-chan struct{} {
	return c.done
}

func (c *trackingContext) Err() error {
	return nil
}

func (c *trackingContext) AfterFunc(fn func()) func() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.pending == nil {
		c.pending = make(map[int]func())
	}
	id := c.nextID
	c.nextID++
	c.pending[id] = fn

	return func() bool {
		c.mu.Lock()
		defer c.mu.Unlock()
		if _, exists := c.pending[id]; !exists {
			return false
		}
		delete(c.pending, id)
		c.stopped++
		return true
	}
}

func (c *trackingContext) stoppedCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.stopped
}

func TestResumeContextPreservesRunValuesAndOverridesWithCommandValues(t *testing.T) {
	type key string

	run := context.WithValue(context.Background(), key("run"), "run-value")
	run = context.WithValue(run, key("shared"), "run-value")
	command := context.WithValue(context.Background(), key("command"), "command-value")
	command = context.WithValue(command, key("shared"), "command-value")

	ctx, cleanup := newResumeContext(run, command)
	defer cleanup()

	if got := ctx.Value(key("run")); got != "run-value" {
		t.Fatalf("expected run value, got %v", got)
	}
	if got := ctx.Value(key("command")); got != "command-value" {
		t.Fatalf("expected command value, got %v", got)
	}
	if got := ctx.Value(key("shared")); got != "command-value" {
		t.Fatalf("expected command value to override run value, got %v", got)
	}
}

func TestResumeContextUsesRunContextWhenCommandIsNil(t *testing.T) {
	type key struct{}

	run := context.WithValue(context.Background(), key{}, "run-value")

	ctx, cleanup := newResumeContext(run, nil)
	defer cleanup()

	if ctx != run {
		t.Fatal("expected nil command context to reuse run context")
	}
}

func TestResumeContextHonorsRunAndCommandCancellation(t *testing.T) {
	for _, tc := range []struct {
		name   string
		cancel func(context.CancelCauseFunc, context.CancelCauseFunc, error)
	}{
		{
			name: "run",
			cancel: func(cancelRun, _ context.CancelCauseFunc, cause error) {
				cancelRun(cause)
			},
		},
		{
			name: "command",
			cancel: func(_, cancelCommand context.CancelCauseFunc, cause error) {
				cancelCommand(cause)
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			run, cancelRun := context.WithCancelCause(context.Background())
			command, cancelCommand := context.WithCancelCause(context.Background())
			ctx, cleanup := newResumeContext(run, command)
			defer cleanup()

			cause := errors.New(tc.name + " canceled")
			tc.cancel(cancelRun, cancelCommand, cause)

			select {
			case <-ctx.Done():
			case <-time.After(time.Second):
				t.Fatal("resume context was not canceled")
			}
			if !errors.Is(ctx.Err(), context.Canceled) {
				t.Fatalf("expected canceled error, got %v", ctx.Err())
			}
			if !errors.Is(context.Cause(ctx), cause) {
				t.Fatalf("expected cancellation cause %v, got %v", cause, context.Cause(ctx))
			}
		})
	}
}

func TestResumeContextReportsEarliestDeadline(t *testing.T) {
	now := time.Now()
	run, cancelRun := context.WithDeadline(context.Background(), now.Add(time.Hour))
	defer cancelRun()
	command, cancelCommand := context.WithDeadline(context.Background(), now.Add(2*time.Hour))
	defer cancelCommand()

	ctx, cleanup := newResumeContext(run, command)
	defer cleanup()

	deadline, ok := ctx.Deadline()
	if !ok || !deadline.Equal(now.Add(time.Hour)) {
		t.Fatalf("expected earliest deadline, got %v, %t", deadline, ok)
	}
}

func TestResumeContextCleanupStopsCancellationHooks(t *testing.T) {
	run := &trackingContext{Context: context.Background(), done: make(chan struct{})}
	command := &trackingContext{Context: context.Background(), done: make(chan struct{})}

	_, cleanup := newResumeContext(run, command)
	cleanup()
	cleanup()

	if got := run.stoppedCount(); got != 1 {
		t.Fatalf("expected one stopped run hook, got %d", got)
	}
	if got := command.stoppedCount(); got != 1 {
		t.Fatalf("expected one stopped command hook, got %d", got)
	}
}
