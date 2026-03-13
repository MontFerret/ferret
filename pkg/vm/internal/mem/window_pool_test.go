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
