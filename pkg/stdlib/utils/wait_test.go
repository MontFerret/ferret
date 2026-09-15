package utils

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestWaitDurationContract(t *testing.T) {
	if value, err := Wait(t.Context(), runtime.ZeroDuration); err != nil || value != runtime.None {
		t.Fatalf("Wait(0s) = %v, %v", value, err)
	}
	if value, err := Wait(t.Context(), runtime.NewInt(0)); err != nil || value != runtime.None {
		t.Fatalf("Wait(0ms) = %v, %v", value, err)
	}
	if value, err := Wait(t.Context(), runtime.NewString("0s")); err != nil || value != runtime.None {
		t.Fatalf("Wait(\"0s\") = %v, %v", value, err)
	}
	if _, err := Wait(t.Context(), runtime.NewDuration(-time.Nanosecond)); !errors.Is(err, runtime.ErrInvalidArgument) {
		t.Fatalf("Wait(negative) error = %v", err)
	}
}

func TestWaitCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	if _, err := Wait(ctx, runtime.NewDuration(time.Hour)); !errors.Is(err, context.Canceled) {
		t.Fatalf("Wait(canceled) error = %v", err)
	}
}

func TestWaitCancellationWhileWaiting(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	observed := &observedWaitContext{Context: ctx, ready: make(chan struct{})}
	finished := make(chan error, 1)
	go func() {
		_, err := Wait(observed, runtime.NewDuration(time.Hour))
		finished <- err
	}()

	select {
	case <-observed.ready:
		cancel()
	case <-time.After(5 * time.Second):
		t.Fatal("Wait never reached its cancellation select")
	}

	select {
	case err := <-finished:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("Wait error = %v, want cancellation", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Wait did not observe cancellation")
	}
}

func TestWaitDeadline(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), time.Millisecond)
	defer cancel()

	if _, err := Wait(ctx, runtime.NewDuration(time.Hour)); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Wait error = %v, want deadline exceeded", err)
	}
}
