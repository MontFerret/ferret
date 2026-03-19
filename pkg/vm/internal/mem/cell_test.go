package mem

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestCellStoreReset_InvalidatesStaleHandleForGet(t *testing.T) {
	var store CellStore

	handle := store.New(runtime.NewInt(1))

	store.Reset()

	if _, ok := store.Get(handle); ok {
		t.Fatal("expected reset to invalidate stale handle lookup")
	}
}

func TestCellStoreReset_InvalidatesStaleHandleForSet(t *testing.T) {
	var store CellStore

	stale := store.New(runtime.NewString("old"))

	store.Reset()

	fresh := store.New(runtime.NewString("fresh"))

	if ok := store.Set(stale, runtime.NewString("stale")); ok {
		t.Fatal("expected reset to invalidate stale handle writes")
	}

	got, ok := store.Get(fresh)
	if !ok {
		t.Fatal("expected fresh handle to remain valid after stale write attempt")
	}

	if got != runtime.NewString("fresh") {
		t.Fatalf("unexpected fresh cell value: got %v, want %v", got, runtime.NewString("fresh"))
	}
}

func TestCellStoreReset_InvalidatesStaleHandleForDelete(t *testing.T) {
	var store CellStore

	stale := store.New(runtime.NewString("old"))

	store.Reset()

	fresh := store.New(runtime.NewString("fresh"))

	if _, ok := store.Delete(stale); ok {
		t.Fatal("expected reset to invalidate stale handle deletes")
	}

	got, ok := store.Get(fresh)
	if !ok {
		t.Fatal("expected fresh handle to remain available after stale delete attempt")
	}

	if got != runtime.NewString("fresh") {
		t.Fatalf("unexpected fresh cell value: got %v, want %v", got, runtime.NewString("fresh"))
	}
}

func TestCellStoreReset_NewHandleUsesNewTokenAndReusedSlot(t *testing.T) {
	var store CellStore

	oldHandle := store.New(runtime.NewString("before-reset"))

	store.Reset()

	newHandle := store.New(runtime.NewString("after-reset"))

	if got, want := newHandle.slot, oldHandle.slot; got != want {
		t.Fatalf("expected slot numbering to restart after reset: got %d, want %d", got, want)
	}

	if got, want := newHandle.generation, oldHandle.generation+1; got != want {
		t.Fatalf("expected generation to advance after reset: got %d, want %d", got, want)
	}

	if got, want := newHandle.ID(), oldHandle.ID(); got == want {
		t.Fatalf("expected canonical handle token to change after reset: got %d, want != %d", got, want)
	}

	got, ok := store.Get(newHandle)
	if !ok {
		t.Fatal("expected new handle to resolve after reset")
	}

	if got != runtime.NewString("after-reset") {
		t.Fatalf("unexpected new cell value: got %v, want %v", got, runtime.NewString("after-reset"))
	}
}
