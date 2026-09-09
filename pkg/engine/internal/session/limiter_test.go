package session

import (
	"context"
	"errors"
	"testing"
	"testing/synctest"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestUnlimitedLimiterStillChecksAdmission(t *testing.T) {
	for _, limiter := range []*Limiter{nil, NewLimiter(0), NewLimiter(-1)} {
		closed := make(chan struct{})
		for range 3 {
			if err := limiter.Acquire(t.Context(), closed); err != nil {
				t.Fatal(err)
			}
		}

		limiter.Release()
		close(closed)
		if err := limiter.Acquire(t.Context(), closed); !errors.Is(err, runtime.ErrInvalidOperation) {
			t.Fatalf("closed Plan admitted unlimited session: %v", err)
		}

		ctx, cancel := context.WithCancel(t.Context())
		cancel()
		if err := limiter.Acquire(ctx, closed); err != context.Canceled {
			t.Fatalf("cancellation lost precedence over Plan closure: %v", err)
		}
	}
}

func TestLimiterWaitObservesReleaseCancellationAndPlanClose(t *testing.T) {
	for _, action := range []string{"release", "cancel", "close"} {
		t.Run(action, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				limiter := NewLimiter(1)
				closed := make(chan struct{})
				if err := limiter.Acquire(t.Context(), closed); err != nil {
					t.Fatal(err)
				}

				ctx, cancel := context.WithCancel(t.Context())
				defer cancel()
				finished := false
				var err error
				go func() {
					err = limiter.Acquire(ctx, closed)
					finished = true
				}()

				synctest.Wait()
				if finished {
					t.Fatal("admission bypassed occupied capacity")
				}

				var want error
				switch action {
				case "release":
					limiter.Release()
				case "cancel":
					cancel()
					want = context.Canceled
				case "close":
					close(closed)
					want = runtime.ErrInvalidOperation
				}

				synctest.Wait()
				if !finished || !errors.Is(err, want) {
					t.Fatalf("wait finished=%t error=%v, want %v", finished, err, want)
				}

				if got := len(limiter.ch); got != 1 {
					t.Fatalf("wait changed another session's permit: %d", got)
				}

				limiter.Release()
				limiter.Release()
			})
		})
	}
}
