package vm

import (
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/internal/debugpoint"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/frame"
)

func TestDebugExecutionFrameLocalsUseCallerRegistersAndCells(t *testing.T) {
	callerPoint := bytecode.DebugPoint{
		ID:         1,
		PC:         10,
		FunctionID: bytecode.NoFunction,
		Bindings: []bytecode.DebugBinding{
			{Name: "caller", Register: bytecode.NewRegister(0)},
			{Name: "captured", Register: bytecode.NewRegister(1), Mutable: true, Cell: true},
		},
	}
	currentPoint := bytecode.DebugPoint{
		ID:         2,
		PC:         20,
		FunctionID: 0,
		Bindings: []bytecode.DebugBinding{
			{Name: "current", Register: bytecode.NewRegister(0)},
		},
	}
	points, err := debugpoint.New([]bytecode.DebugPoint{callerPoint, currentPoint})
	if err != nil {
		t.Fatal(err)
	}

	instance := &VM{program: &bytecode.Program{}}
	handle := instance.state.cells.New(runtime.NewInt(42))
	instance.state.registers = []runtime.Value{runtime.NewInt(7)}
	instance.state.frames.Push(frame.CallFrame{
		FnID:            0,
		ReturnPC:        12,
		HasCallSite:     true,
		CallerRegisters: []runtime.Value{runtime.NewInt(3), handle},
	})
	execution := &debugExecution{
		vm:      instance,
		points:  points,
		current: &currentPoint,
		status:  DebugExecutionPaused,
	}

	locals, err := execution.FrameLocals(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(locals) != 2 || locals[0].Name != "caller" || locals[0].Value != runtime.NewInt(3) ||
		locals[1].Name != "captured" || !locals[1].Mutable || locals[1].Value != runtime.NewInt(42) {
		t.Fatalf("unexpected caller locals: %#v", locals)
	}

	if !instance.state.cells.Set(handle, runtime.NewInt(43)) {
		t.Fatal("failed to update captured cell")
	}
	locals, err = execution.FrameLocals(1)
	if err != nil {
		t.Fatal(err)
	}
	if locals[1].Value != runtime.NewInt(43) {
		t.Fatalf("caller frame did not observe cell update: %#v", locals)
	}

	if _, err := execution.FrameLocals(-1); !errors.Is(err, runtime.ErrInvalidArgument) {
		t.Fatalf("expected negative frame rejection, got %v", err)
	}
	if _, err := execution.FrameLocals(2); !errors.Is(err, runtime.ErrNotFound) {
		t.Fatalf("expected missing frame rejection, got %v", err)
	}

	execution.status = DebugExecutionRunning
	if _, err := execution.FrameLocals(0); !errors.Is(err, runtime.ErrInvalidOperation) {
		t.Fatalf("expected running execution rejection, got %v", err)
	}
}
