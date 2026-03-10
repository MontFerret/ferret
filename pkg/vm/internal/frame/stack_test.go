package frame

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestCallStackUnwindToProtectedFrame_ReclaimsRegisters(t *testing.T) {
	var stack CallStack
	stack.Init(6)

	lowerRegs := make([]runtime.Value, 2)
	protectedRegs := make([]runtime.Value, 3)
	aboveRegs1 := make([]runtime.Value, 4)
	aboveRegs2 := make([]runtime.Value, 5)
	activeRegs := make([]runtime.Value, 6)

	protectedRegs[1] = runtime.True

	stack.Push(CallFrame{
		ReturnPC:   10,
		ReturnDest: bytecode.NewRegister(0),
		Registers:  lowerRegs,
		Protected:  false,
	})
	stack.Push(CallFrame{
		ReturnPC:   20,
		ReturnDest: bytecode.NewRegister(1),
		Registers:  protectedRegs,
		Protected:  true,
	})
	stack.Push(CallFrame{
		ReturnPC:   30,
		ReturnDest: bytecode.NewRegister(0),
		Registers:  aboveRegs1,
		Protected:  false,
	})
	stack.Push(CallFrame{
		ReturnPC:   40,
		ReturnDest: bytecode.NewRegister(0),
		Registers:  aboveRegs2,
		Protected:  false,
	})

	registers, pc, ok := stack.UnwindToProtectedFrame(activeRegs)
	if !ok {
		t.Fatal("expected protected unwind to succeed")
	}

	if got, want := pc, 20; got != want {
		t.Fatalf("unexpected pc after unwind: got %d, want %d", got, want)
	}

	if got, want := stack.Len(), 1; got != want {
		t.Fatalf("unexpected frame depth after unwind: got %d, want %d", got, want)
	}

	remaining := stack.Top()
	if remaining == nil {
		t.Fatal("expected remaining frame after unwind")
	}

	if got, want := remaining.ReturnPC, 10; got != want {
		t.Fatalf("unexpected surviving frame returnPC: got %d, want %d", got, want)
	}

	if got, want := registers[1], runtime.None; got != want {
		t.Fatalf("expected protected return destination to be reset, got %v", got)
	}

	reused4 := stack.AcquireRegisters(4)
	if len(reused4) != 4 {
		t.Fatalf("unexpected pooled registers length: got %d, want %d", len(reused4), 4)
	}
	if &reused4[0] != &aboveRegs1[0] {
		t.Fatal("expected frame registers of size 4 to be reclaimed")
	}

	reused5 := stack.AcquireRegisters(5)
	if len(reused5) != 5 {
		t.Fatalf("unexpected pooled registers length: got %d, want %d", len(reused5), 5)
	}
	if &reused5[0] != &aboveRegs2[0] {
		t.Fatal("expected frame registers of size 5 to be reclaimed")
	}

	reused6 := stack.AcquireRegisters(6)
	if len(reused6) != 6 {
		t.Fatalf("unexpected pooled registers length: got %d, want %d", len(reused6), 6)
	}
	if &reused6[0] != &activeRegs[0] {
		t.Fatal("expected active registers of size 6 to be reclaimed")
	}
}

func TestCallStackReturnToCaller_ReusesRegisters(t *testing.T) {
	var stack CallStack
	stack.Init(3)

	callerRegs := make([]runtime.Value, 2)
	activeRegs := make([]runtime.Value, 3)

	stack.Push(CallFrame{
		ReturnPC:   7,
		ReturnDest: bytecode.NewRegister(1),
		Registers:  callerRegs,
		Protected:  false,
	})

	retVal := runtime.True
	registers, pc, ok := stack.ReturnToCaller(activeRegs, retVal)
	if !ok {
		t.Fatal("expected return to succeed")
	}

	if got, want := pc, 7; got != want {
		t.Fatalf("unexpected pc after return: got %d, want %d", got, want)
	}

	if got, want := registers[1], retVal; got != want {
		t.Fatalf("unexpected return destination: got %v, want %v", got, want)
	}

	reused := stack.AcquireRegisters(3)
	if len(reused) != 3 {
		t.Fatalf("unexpected pooled registers length: got %d, want %d", len(reused), 3)
	}
	if &reused[0] != &activeRegs[0] {
		t.Fatal("expected active registers to be reclaimed")
	}
}

func TestCallStackUnwindToProtectedFrame_NoProtectedFrames(t *testing.T) {
	var stack CallStack
	stack.Init(4)

	activeRegs := make([]runtime.Value, 4)
	stack.Push(CallFrame{
		ReturnPC:   10,
		ReturnDest: bytecode.NewRegister(0),
		Registers:  make([]runtime.Value, 2),
		Protected:  false,
	})

	if _, _, ok := stack.UnwindToProtectedFrame(activeRegs); ok {
		t.Fatal("expected protected unwind to fail without protected frames")
	}

	if got, want := stack.Len(), 1; got != want {
		t.Fatalf("unexpected frame depth after failed unwind: got %d, want %d", got, want)
	}

	reused := stack.AcquireRegisters(4)
	if len(reused) != 4 {
		t.Fatalf("unexpected pooled registers length: got %d, want %d", len(reused), 4)
	}
	if &reused[0] == &activeRegs[0] {
		t.Fatal("did not expect active registers to be released on failed unwind")
	}
}

func TestCallStackReturnToCaller_EmptyStackDoesNotRelease(t *testing.T) {
	var stack CallStack
	stack.Init(2)

	activeRegs := make([]runtime.Value, 2)
	if _, _, ok := stack.ReturnToCaller(activeRegs, runtime.True); ok {
		t.Fatal("expected return to fail on empty stack")
	}

	reused := stack.AcquireRegisters(2)
	if len(reused) != 2 {
		t.Fatalf("unexpected pooled registers length: got %d, want %d", len(reused), 2)
	}
	if &reused[0] == &activeRegs[0] {
		t.Fatal("did not expect active registers to be released on empty return")
	}
}

func TestCallStackSetTopFnID(t *testing.T) {
	var stack CallStack
	stack.Init(1)

	if ok := stack.SetTopFnID(10); ok {
		t.Fatal("expected SetTopFnID to fail on empty stack")
	}

	stack.Push(CallFrame{
		ReturnPC:   5,
		ReturnDest: bytecode.NewRegister(0),
		Registers:  make([]runtime.Value, 1),
		Protected:  false,
		FnID:       1,
	})

	if ok := stack.SetTopFnID(42); !ok {
		t.Fatal("expected SetTopFnID to succeed with a frame present")
	}

	top := stack.Top()
	if top == nil {
		t.Fatal("expected top frame to be present")
	}
	if got, want := top.FnID, 42; got != want {
		t.Fatalf("unexpected fnID after update: got %d, want %d", got, want)
	}
}

func TestCallStackTraceEntries_OrderAndMetadata(t *testing.T) {
	var stack CallStack
	stack.Init(1)

	stack.Push(CallFrame{
		FnID:        1,
		FnName:      "outer",
		CallSitePC:  11,
		HasCallSite: true,
	})
	stack.Push(CallFrame{
		FnID:        2,
		FnName:      "middle",
		CallSitePC:  22,
		HasCallSite: true,
	})
	stack.Push(CallFrame{
		FnID:        3,
		FnName:      "inner",
		CallSitePC:  33,
		HasCallSite: true,
	})

	traces := stack.TraceEntries()
	if got, want := len(traces), 3; got != want {
		t.Fatalf("unexpected trace count: got %d, want %d", got, want)
	}

	if got, want := traces[0].FnName, "inner"; got != want {
		t.Fatalf("unexpected nearest trace name: got %q, want %q", got, want)
	}

	if got, want := traces[2].FnName, "outer"; got != want {
		t.Fatalf("unexpected farthest trace name: got %q, want %q", got, want)
	}

	pcs := stack.CallSitePCs()
	if got, want := len(pcs), 3; got != want {
		t.Fatalf("unexpected pc count: got %d, want %d", got, want)
	}

	if got, want := pcs[0], 33; got != want {
		t.Fatalf("unexpected nearest pc: got %d, want %d", got, want)
	}
}
