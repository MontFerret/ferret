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
	if _, err := Wait(t.Context(), runtime.NewInt(0)); !errors.Is(err, runtime.ErrInvalidType) {
		t.Fatalf("Wait(Int) error = %v", err)
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
