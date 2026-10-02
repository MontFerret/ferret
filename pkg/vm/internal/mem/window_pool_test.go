package mem

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestWindowPoolAcquireNewWindowInitializesSlotsWithNone(t *testing.T) {
	pool := NewWindowPool(4)

	reg := pool.Acquire(3)
	if got, want := len(reg), 3; got != want {
		t.Fatalf("unexpected register window size: got %d, want %d", got, want)
	}

	for i := range reg {
		if got := reg[i]; got != runtime.None {
			t.Fatalf("expected slot %d to be runtime.None, got %v", i, got)
		}
	}
}

func TestWindowPoolReleaseAcquireReusesNormalizedWindow(t *testing.T) {
	pool := NewWindowPool(4)

	reg := make([]runtime.Value, 3)
	reg[0] = runtime.NewInt(1)
	reg[1] = runtime.NewInt(2)
	reg[2] = runtime.NewInt(3)

	pool.Release(reg)
	reused := pool.Acquire(3)

	if &reused[0] != &reg[0] {
		t.Fatal("expected released register window to be reused")
	}

	for i := range reused {
		if got := reused[i]; got != runtime.None {
			t.Fatalf("expected reused slot %d to be runtime.None, got %v", i, got)
		}
	}
}

func TestWindowPoolZeroAndOversizeBehavior(t *testing.T) {
	pool := NewWindowPool(2)

	if got := pool.Acquire(0); got != nil {
		t.Fatalf("expected nil window for zero size, got %#v", got)
	}

	oversize := pool.Acquire(3)
	if got, want := len(oversize), 3; got != want {
		t.Fatalf("unexpected oversize window length: got %d, want %d", got, want)
	}

	pool.Release(oversize)
	reused := pool.Acquire(3)
	if &reused[0] == &oversize[0] {
		t.Fatal("did not expect oversize window to be pooled")
	}

	pool.Release(oversize[:1])
	small := pool.Acquire(1)
	if &small[0] == &oversize[0] {
		t.Fatal("shortening an oversize window must not make it poolable")
	}
}

func TestWindowPoolRelease_ScrubsWithoutClosingValues(t *testing.T) {
	pool := NewWindowPool(2)
	closer := newTestCloser("window")
	reg := []runtime.Value{closer}

	pool.Release(reg)

	if got := closer.closed; got != 0 {
		t.Fatalf("expected pooled window release to avoid closing values, got %d closes", got)
	}

	reused := pool.Acquire(1)
	if got := reused[0]; got != runtime.None {
		t.Fatalf("expected pooled slot to be scrubbed to runtime.None, got %v", got)
	}
}

func TestWindowPoolReleaseShrunkWindowReusesBackingCapacity(t *testing.T) {
	pool := NewWindowPool(8)
	reg := pool.Acquire(8)
	closer := newTestCloser("discarded tail slot")
	reg[7] = closer

	// Tail replacement reuses a larger window with the callee's shorter length.
	pool.Release(reg[:4])
	reused := pool.Acquire(8)
	if len(reused) != 8 {
		t.Fatalf("reused window length = %d, want 8", len(reused))
	}
	if &reused[0] != &reg[0] {
		t.Fatal("shrinking a window must preserve reuse of its backing capacity")
	}

	for i, value := range reused {
		if value != runtime.None {
			t.Fatalf("released backing slot %d = %v, want runtime.None", i, value)
		}
	}

	if closer.closed != 0 {
		t.Fatalf("window storage closed a resource %d times", closer.closed)
	}
}

func TestWindowPoolRepeatedShrinkingKeepsRetainedWindowsBounded(t *testing.T) {
	pool := NewWindowPool(8)
	for range 64 {
		reg := pool.Acquire(8)
		pool.Release(reg[:4])
	}

	retained := 0
	for _, bucket := range pool.buckets {
		retained += len(bucket)
	}

	if retained != 1 {
		t.Fatalf("repeatedly shrinking one active window retained %d windows, want 1", retained)
	}
}
