package core

import (
	"errors"
	"strings"
	"testing"
)

func TestEmitterPatchxPanicsOnUnmarkedLabel(t *testing.T) {
	emitter := NewEmitter()
	label := emitter.NewLabel("missing")

	assertInvariantPanic(t, func() {
		emitter.Patchx(label, 1)
	}, "label not marked")
}

func TestLoopTableRequiredParentPanicsWhenParentMissing(t *testing.T) {
	table := NewLoopTable(NewRegisterAllocator())
	table.Push(&Loop{Allocate: false})
	table.Push(&Loop{Allocate: false})

	assertInvariantPanic(t, func() {
		table.RequiredParent(table.Depth())
	}, "parent loop not found in loop table")
}

func assertInvariantPanic(t *testing.T, fn func(), want string) {
	t.Helper()

	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatal("expected invariant panic")
		}

		err, ok := recovered.(error)
		if !ok {
			t.Fatalf("expected panic error, got %T", recovered)
		}

		var invariant *InvariantViolation
		if !errors.As(err, &invariant) {
			t.Fatalf("expected invariant violation, got %T", err)
		}

		if !strings.Contains(err.Error(), want) {
			t.Fatalf("unexpected invariant panic: got %q want substring %q", err.Error(), want)
		}
	}()

	fn()
}
