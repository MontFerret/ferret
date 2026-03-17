package mem

import (
	"testing"
)

func TestAliasTracker_IncDec(t *testing.T) {
	closer := newTestCloser("a")
	var tracker AliasTracker

	tracker.Inc(closer)
	if got := tracker.Count(closer); got != 1 {
		t.Fatalf("expected count 1 after Inc, got %d", got)
	}

	tracker.Inc(closer)
	if got := tracker.Count(closer); got != 2 {
		t.Fatalf("expected count 2 after second Inc, got %d", got)
	}

	remaining := tracker.Dec(closer)
	if remaining != 1 {
		t.Fatalf("expected Dec to return 1, got %d", remaining)
	}

	remaining = tracker.Dec(closer)
	if remaining != 0 {
		t.Fatalf("expected Dec to return 0, got %d", remaining)
	}

	if got := tracker.Count(closer); got != 0 {
		t.Fatalf("expected count 0 after decrementing to zero, got %d", got)
	}
}

func TestAliasTracker_DecUntracked(t *testing.T) {
	closer := newTestCloser("untracked")
	var tracker AliasTracker

	if got := tracker.Dec(closer); got != 0 {
		t.Fatalf("expected Dec on untracked closer to return 0, got %d", got)
	}
}

func TestAliasTracker_MultipleClosers(t *testing.T) {
	a := newTestCloser("a")
	b := newTestCloser("b")
	var tracker AliasTracker

	tracker.Inc(a)
	tracker.Inc(a)
	tracker.Inc(b)

	if got := tracker.Count(a); got != 2 {
		t.Fatalf("expected a count 2, got %d", got)
	}

	if got := tracker.Count(b); got != 1 {
		t.Fatalf("expected b count 1, got %d", got)
	}

	tracker.Dec(a)
	if got := tracker.Count(a); got != 1 {
		t.Fatalf("expected a count 1 after Dec, got %d", got)
	}

	if got := tracker.Count(b); got != 1 {
		t.Fatalf("expected b count unchanged at 1, got %d", got)
	}
}

func TestAliasTracker_Reset(t *testing.T) {
	a := newTestCloser("a")
	var tracker AliasTracker

	tracker.Inc(a)
	tracker.Inc(a)
	tracker.Reset()

	if got := tracker.Count(a); got != 0 {
		t.Fatalf("expected count 0 after Reset, got %d", got)
	}
}

func TestAliasTracker_ZeroValue(t *testing.T) {
	a := newTestCloser("a")
	var tracker AliasTracker

	// Zero-value tracker should handle all operations safely
	if got := tracker.Count(a); got != 0 {
		t.Fatalf("expected count 0 on zero-value tracker, got %d", got)
	}

	if got := tracker.Dec(a); got != 0 {
		t.Fatalf("expected Dec 0 on zero-value tracker, got %d", got)
	}

	tracker.Reset() // should not panic
}
