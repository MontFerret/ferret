package vm

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/mem"
)

func TestLifecycle_CellOverwriteAndCleanupCloseValuesExactlyOnce(t *testing.T) {
	instance := mustNewVM(t, &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  3,
	})

	state := mustAcquireRunState(t, instance)

	original := newTrackingCloser("original")
	replacement := newTrackingCloser("replacement")

	state.writeProducedRegister(bytecode.NewRegister(0), original)

	handleValue := state.makeCell(bytecode.NewRegister(1), state.registers[bytecode.NewRegister(0)])
	handle, ok := handleValue.(mem.CellHandle)
	if !ok {
		t.Fatalf("expected cell handle, got %T", handleValue)
	}

	state.clearRegister(bytecode.NewRegister(0))

	state.writeProducedRegister(bytecode.NewRegister(2), replacement)

	if err := state.storeCell(handle, state.registers[bytecode.NewRegister(2)]); err != nil {
		t.Fatalf("unexpected storeCell error: %v", err)
	}

	if got, want := countDeferredClosers(&state.deferred), 1; got != want {
		t.Fatalf("expected deferred closers after overwriting cell value: got %d, want %d", got, want)
	}

	if state.owned.Owns(original) {
		t.Fatal("expected original cell value ownership to be released after overwrite")
	}

	if got := original.closed; got != 0 {
		t.Fatalf("expected original to stay deferred until run end, got %d closes", got)
	}

	state.clearRegister(bytecode.NewRegister(2))

	if !state.owned.Owns(replacement) {
		t.Fatal("expected replacement to stay owned while still referenced by the cell")
	}

	if got, want := countDeferredClosers(&state.deferred), 1; got != want {
		t.Fatalf("expected deferred closers to stay unchanged before cell cleanup: got %d, want %d", got, want)
	}

	state.cleanupCurrentCells()

	if state.owned.Owns(replacement) {
		t.Fatal("expected replacement ownership to end after cell cleanup")
	}

	if got, want := countDeferredClosers(&state.deferred), 2; got != want {
		t.Fatalf("expected both cell values to be deferred after cleanup: got %d, want %d", got, want)
	}

	if got := replacement.closed; got != 0 {
		t.Fatalf("expected replacement to stay deferred until run end, got %d closes", got)
	}

	state.endRun()

	if got := original.closed; got != 1 {
		t.Fatalf("expected original to close exactly once, got %d closes", got)
	}

	if got := replacement.closed; got != 1 {
		t.Fatalf("expected replacement to close exactly once, got %d closes", got)
	}
}
