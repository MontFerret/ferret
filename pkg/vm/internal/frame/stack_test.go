package frame

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestCallStackPushPopTopAndLen(t *testing.T) {
	var stack CallStack

	lowerRegs := make([]runtime.Value, 2)
	upperRegs := make([]runtime.Value, 3)

	stack.Push(CallFrame{
		FnName:          "outer",
		CallerRegisters: lowerRegs,
		ReturnPC:        10,
		ReturnDest:      bytecode.NewRegister(1),
	})
	stack.Push(CallFrame{
		FnName:           "inner",
		CallerRegisters:  upperRegs,
		ReturnPC:         20,
		ReturnDest:       bytecode.NewRegister(2),
		RecoveryBoundary: true,
	})

	if got, want := stack.Len(), 2; got != want {
		t.Fatalf("unexpected stack length: got %d, want %d", got, want)
	}

	top := stack.Top()
	if top == nil {
		t.Fatal("expected top frame")
	}
	if got, want := top.FnName, "inner"; got != want {
		t.Fatalf("unexpected top frame name: got %q, want %q", got, want)
	}

	backing := stack.frames[:len(stack.frames)]
	popped, ok := stack.Pop()
	if !ok {
		t.Fatal("expected pop to succeed")
	}
	if got, want := popped.FnName, "inner"; got != want {
		t.Fatalf("unexpected popped frame name: got %q, want %q", got, want)
	}
	if got, want := stack.Len(), 1; got != want {
		t.Fatalf("unexpected stack length after pop: got %d, want %d", got, want)
	}
	if backing[1].CallerRegisters != nil || backing[1].FnName != "" || backing[1].FnID != 0 || backing[1].CallSitePC != 0 || backing[1].ReturnPC != 0 || backing[1].ReturnDest != 0 || backing[1].HasCallSite || backing[1].RecoveryBoundary {
		t.Fatalf("expected popped backing slot to be zeroed, got %+v", backing[1])
	}

	popped, ok = stack.Pop()
	if !ok {
		t.Fatal("expected second pop to succeed")
	}
	if got, want := popped.FnName, "outer"; got != want {
		t.Fatalf("unexpected second popped frame name: got %q, want %q", got, want)
	}
	if got := stack.Top(); got != nil {
		t.Fatalf("expected empty stack after pops, got %+v", got)
	}
	if got, want := stack.Len(), 0; got != want {
		t.Fatalf("unexpected empty stack length: got %d, want %d", got, want)
	}
}

func TestCallStackNearestRecoveryBoundary(t *testing.T) {
	var stack CallStack

	if got, want := stack.NearestRecoveryBoundary(), -1; got != want {
		t.Fatalf("unexpected empty boundary index: got %d, want %d", got, want)
	}

	stack.Push(CallFrame{FnName: "outer"})
	stack.Push(CallFrame{FnName: "guard", RecoveryBoundary: true})
	stack.Push(CallFrame{FnName: "inner"})
	stack.Push(CallFrame{FnName: "guard2", RecoveryBoundary: true})

	if got, want := stack.NearestRecoveryBoundary(), 3; got != want {
		t.Fatalf("unexpected topmost boundary index: got %d, want %d", got, want)
	}

	if _, ok := stack.Pop(); !ok {
		t.Fatal("expected pop to succeed")
	}
	if got, want := stack.NearestRecoveryBoundary(), 1; got != want {
		t.Fatalf("unexpected boundary index after pop: got %d, want %d", got, want)
	}
}

func TestCallStackSetTopMetadata(t *testing.T) {
	var stack CallStack

	if ok := stack.SetTopFnID(10); ok {
		t.Fatal("expected SetTopFnID to fail on empty stack")
	}
	if ok := stack.SetTopCall(10, "fn", 12); ok {
		t.Fatal("expected SetTopCall to fail on empty stack")
	}

	stack.Push(CallFrame{
		FnID:            1,
		FnName:          "outer",
		CallerRegisters: make([]runtime.Value, 1),
	})

	if ok := stack.SetTopFnID(42); !ok {
		t.Fatal("expected SetTopFnID to succeed")
	}
	if ok := stack.SetTopCall(42, "renamed", 33); !ok {
		t.Fatal("expected SetTopCall to succeed")
	}

	top := stack.Top()
	if top == nil {
		t.Fatal("expected top frame")
	}
	if got, want := top.FnID, 42; got != want {
		t.Fatalf("unexpected fn id: got %d, want %d", got, want)
	}
	if got, want := top.FnName, "renamed"; got != want {
		t.Fatalf("unexpected fn name: got %q, want %q", got, want)
	}
	if got, want := top.CallSitePC, 33; got != want {
		t.Fatalf("unexpected callsite pc: got %d, want %d", got, want)
	}
	if !top.HasCallSite {
		t.Fatal("expected top frame callsite metadata to be marked present")
	}
}

func TestCallStackTraceEntriesOrderAndMetadata(t *testing.T) {
	var stack CallStack

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
		t.Fatalf("unexpected nearest callsite pc: got %d, want %d", got, want)
	}
	if got, want := pcs[2], 11; got != want {
		t.Fatalf("unexpected farthest callsite pc: got %d, want %d", got, want)
	}
}
