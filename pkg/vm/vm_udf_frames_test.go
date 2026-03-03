package vm

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/frame"
)

func TestUnwindToProtected_ReclaimsDiscardedFrameRegisters(t *testing.T) {
	instance := New(&bytecode.Program{
		Registers: 1,
		Functions: bytecode.Functions{
			UserDefined: []bytecode.UDF{
				{Registers: 6},
			},
		},
	})

	lowerRegs := make([]runtime.Value, 2)
	protectedRegs := make([]runtime.Value, 3)
	aboveRegs1 := make([]runtime.Value, 4)
	aboveRegs2 := make([]runtime.Value, 5)
	activeRegs := make([]runtime.Value, 6)

	protectedRegs[1] = runtime.True

	instance.registers.Values = activeRegs
	instance.frames.Push(frame.CallFrame{
		ReturnPC:   10,
		ReturnDest: bytecode.NewRegister(0),
		Registers:  lowerRegs,
		Protected:  false,
	})
	instance.frames.Push(frame.CallFrame{
		ReturnPC:   20,
		ReturnDest: bytecode.NewRegister(1),
		Registers:  protectedRegs,
		Protected:  true,
	})
	instance.frames.Push(frame.CallFrame{
		ReturnPC:   30,
		ReturnDest: bytecode.NewRegister(0),
		Registers:  aboveRegs1,
		Protected:  false,
	})
	instance.frames.Push(frame.CallFrame{
		ReturnPC:   40,
		ReturnDest: bytecode.NewRegister(0),
		Registers:  aboveRegs2,
		Protected:  false,
	})

	if ok := instance.unwindToProtected(); !ok {
		t.Fatal("expected protected unwind to succeed")
	}

	if got, want := instance.pc, 20; got != want {
		t.Fatalf("unexpected pc after unwind: got %d, want %d", got, want)
	}

	if got, want := instance.frames.Len(), 1; got != want {
		t.Fatalf("unexpected frame depth after unwind: got %d, want %d", got, want)
	}

	remaining := instance.frames.Top()
	if remaining == nil {
		t.Fatal("expected remaining frame after unwind")
	}

	if got, want := remaining.ReturnPC, 10; got != want {
		t.Fatalf("unexpected surviving frame returnPC: got %d, want %d", got, want)
	}

	if got, want := instance.registers.Values[1], runtime.None; got != want {
		t.Fatalf("expected protected return destination to be reset, got %v", got)
	}

	reused4 := instance.frames.AcquireRegisters(4)
	if len(reused4) != 4 {
		t.Fatalf("unexpected pooled registers length: got %d, want %d", len(reused4), 4)
	}
	if &reused4[0] != &aboveRegs1[0] {
		t.Fatal("expected frame registers of size 4 to be reclaimed")
	}

	reused5 := instance.frames.AcquireRegisters(5)
	if len(reused5) != 5 {
		t.Fatalf("unexpected pooled registers length: got %d, want %d", len(reused5), 5)
	}
	if &reused5[0] != &aboveRegs2[0] {
		t.Fatal("expected frame registers of size 5 to be reclaimed")
	}

	reused6 := instance.frames.AcquireRegisters(6)
	if len(reused6) != 6 {
		t.Fatalf("unexpected pooled registers length: got %d, want %d", len(reused6), 6)
	}
	if &reused6[0] != &activeRegs[0] {
		t.Fatal("expected active registers of size 6 to be reclaimed")
	}
}
