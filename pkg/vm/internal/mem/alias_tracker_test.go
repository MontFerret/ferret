package mem

import (
	"testing"
)

func TestAliasTracker_IncDec(t *testing.T) {
	closer := newTestCloser("a")
	key, _, _ := ResourceKeyOf(closer)
	var tracker AliasTracker

	tracker.Inc(key)
	if got := tracker.Count(key); got != 1 {
		t.Fatalf("expected count 1 after Inc, got %d", got)
	}

	tracker.Inc(key)
	if got := tracker.Count(key); got != 2 {
		t.Fatalf("expected count 2 after second Inc, got %d", got)
	}

	remaining := tracker.Dec(key)
	if remaining != 1 {
		t.Fatalf("expected Dec to return 1, got %d", remaining)
	}

	remaining = tracker.Dec(key)
	if remaining != 0 {
		t.Fatalf("expected Dec to return 0, got %d", remaining)
	}

	if got := tracker.Count(key); got != 0 {
		t.Fatalf("expected count 0 after decrementing to zero, got %d", got)
	}
}

func TestAliasTracker_DecUntracked(t *testing.T) {
	closer := newTestCloser("untracked")
	key, _, _ := ResourceKeyOf(closer)
	var tracker AliasTracker

	if got := tracker.Dec(key); got != 0 {
		t.Fatalf("expected Dec on untracked closer to return 0, got %d", got)
	}
}

func TestAliasTracker_MultipleClosers(t *testing.T) {
	a := newTestCloser("a")
	b := newTestCloser("b")
	keyA, _, _ := ResourceKeyOf(a)
	keyB, _, _ := ResourceKeyOf(b)
	var tracker AliasTracker

	tracker.Inc(keyA)
	tracker.Inc(keyA)
	tracker.Inc(keyB)

	if got := tracker.Count(keyA); got != 2 {
		t.Fatalf("expected a count 2, got %d", got)
	}

	if got := tracker.Count(keyB); got != 1 {
		t.Fatalf("expected b count 1, got %d", got)
	}

	tracker.Dec(keyA)
	if got := tracker.Count(keyA); got != 1 {
		t.Fatalf("expected a count 1 after Dec, got %d", got)
	}

	if got := tracker.Count(keyB); got != 1 {
		t.Fatalf("expected b count unchanged at 1, got %d", got)
	}
}

func TestAliasTracker_Reset(t *testing.T) {
	a := newTestCloser("a")
	keyA, _, _ := ResourceKeyOf(a)
	var tracker AliasTracker

	tracker.Inc(keyA)
	tracker.Inc(keyA)
	tracker.Reset()

	if got := tracker.Count(keyA); got != 0 {
		t.Fatalf("expected count 0 after Reset, got %d", got)
	}
}

func TestAliasTracker_Delete(t *testing.T) {
	a := newTestCloser("a")
	b := newTestCloser("b")
	keyA, _, _ := ResourceKeyOf(a)
	keyB, _, _ := ResourceKeyOf(b)
	var tracker AliasTracker

	tracker.Inc(keyA)
	tracker.Inc(keyA)
	tracker.Inc(keyB)
	tracker.Delete(keyA)

	if got := tracker.Count(keyA); got != 0 {
		t.Fatalf("expected deleted key count 0, got %d", got)
	}

	if got := tracker.Count(keyB); got != 1 {
		t.Fatalf("expected unrelated key count 1, got %d", got)
	}
}

func TestAliasTracker_ZeroValue(t *testing.T) {
	a := newTestCloser("a")
	keyA, _, _ := ResourceKeyOf(a)
	var tracker AliasTracker

	// Zero-value tracker should handle all operations safely
	if got := tracker.Count(keyA); got != 0 {
		t.Fatalf("expected count 0 on zero-value tracker, got %d", got)
	}

	if got := tracker.Dec(keyA); got != 0 {
		t.Fatalf("expected Dec 0 on zero-value tracker, got %d", got)
	}

	tracker.Reset() // should not panic
}
