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

func TestLifecycle_MultipleCellsCleanupIndependently(t *testing.T) {
	instance := mustNewVM(t, &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  6,
	})

	state := mustAcquireRunState(t, instance)

	originalA := newTrackingCloser("original-a")
	originalB := newTrackingCloser("original-b")
	replacementA := newTrackingCloser("replacement-a")
	replacementB := newTrackingCloser("replacement-b")

	state.writeProducedRegister(bytecode.NewRegister(0), originalA)

	handleValueA := state.makeCell(bytecode.NewRegister(1), state.registers[bytecode.NewRegister(0)])
	handleA, ok := handleValueA.(mem.CellHandle)
	if !ok {
		t.Fatalf("expected first cell handle, got %T", handleValueA)
	}

	state.writeProducedRegister(bytecode.NewRegister(2), originalB)

	handleValueB := state.makeCell(bytecode.NewRegister(3), state.registers[bytecode.NewRegister(2)])
	handleB, ok := handleValueB.(mem.CellHandle)
	if !ok {
		t.Fatalf("expected second cell handle, got %T", handleValueB)
	}

	state.clearRegister(bytecode.NewRegister(0))
	state.clearRegister(bytecode.NewRegister(2))

	state.writeProducedRegister(bytecode.NewRegister(4), replacementA)

	if err := state.storeCell(handleA, state.registers[bytecode.NewRegister(4)]); err != nil {
		t.Fatalf("unexpected first storeCell error: %v", err)
	}

	if got, want := countDeferredClosers(&state.deferred), 1; got != want {
		t.Fatalf("expected deferred closers after first overwrite: got %d, want %d", got, want)
	}

	state.clearRegister(bytecode.NewRegister(4))

	state.writeProducedRegister(bytecode.NewRegister(5), replacementB)

	if err := state.storeCell(handleB, state.registers[bytecode.NewRegister(5)]); err != nil {
		t.Fatalf("unexpected second storeCell error: %v", err)
	}

	if got, want := countDeferredClosers(&state.deferred), 2; got != want {
		t.Fatalf("expected deferred closers after second overwrite: got %d, want %d", got, want)
	}

	state.clearRegister(bytecode.NewRegister(5))

	if state.owned.Owns(originalA) || state.owned.Owns(originalB) {
		t.Fatal("expected original cell values ownership to be released after overwrite")
	}

	if !state.owned.Owns(replacementA) || !state.owned.Owns(replacementB) {
		t.Fatal("expected replacement values to remain owned while referenced by cells")
	}

	state.cleanupCurrentCells()

	if state.owned.Owns(replacementA) || state.owned.Owns(replacementB) {
		t.Fatal("expected replacement ownership to end after cell cleanup")
	}

	if got, want := countDeferredClosers(&state.deferred), 4; got != want {
		t.Fatalf("expected all cell values to be deferred after cleanup: got %d, want %d", got, want)
	}

	state.endRun()

	if got := originalA.closed; got != 1 {
		t.Fatalf("expected originalA to close exactly once, got %d closes", got)
	}

	if got := originalB.closed; got != 1 {
		t.Fatalf("expected originalB to close exactly once, got %d closes", got)
	}

	if got := replacementA.closed; got != 1 {
		t.Fatalf("expected replacementA to close exactly once, got %d closes", got)
	}

	if got := replacementB.closed; got != 1 {
		t.Fatalf("expected replacementB to close exactly once, got %d closes", got)
	}
}
