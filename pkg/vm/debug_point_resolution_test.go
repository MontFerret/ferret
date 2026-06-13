package vm

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/internal/debugpoint"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/frame"
)

func TestDebugExecutionErrorPointDoesNotCrossFunctionBoundaries(t *testing.T) {
	points := []bytecode.DebugPoint{
		{ID: 4, PC: 2, FunctionID: 0},
		{ID: 7, PC: 4, FunctionID: 1},
		{ID: 2, PC: 6, FunctionID: 0},
	}
	instance := &VM{program: &bytecode.Program{}}
	instance.state.lastPC = 5
	instance.state.frames.Push(frame.CallFrame{FnID: 0})
	execution := &debugExecution{
		vm:     instance,
		points: debugpoint.New(points),
	}

	if got := execution.errorPoint(); got != &points[0] {
		t.Fatalf("resolved runtime error through another function's point: %#v", got)
	}
}
