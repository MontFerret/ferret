package vm

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestElapsedOpcodeReturnsMonotonicDuration(t *testing.T) {
	program := newTestProgram(
		3,
		nil,
		bytecode.NewInstruction(bytecode.OpElapsed, bytecode.NewRegister(0)),
		bytecode.NewInstruction(bytecode.OpElapsed, bytecode.NewRegister(1)),
		bytecode.NewInstruction(bytecode.OpSub, bytecode.NewRegister(2), bytecode.NewRegister(1), bytecode.NewRegister(0)),
		bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(2)),
	)
	instance := mustNewVM(t, program)
	t.Cleanup(func() { _ = instance.Close() })

	result := mustRunResult(t, instance, NewDefaultEnvironment())
	value := mustResultRootAndClose(t, result)
	elapsed, ok := value.(runtime.Duration)
	if !ok {
		t.Fatalf("expected Duration result, got %T", value)
	}
	if elapsed < 0 {
		t.Fatalf("expected monotonic elapsed time, got %s", elapsed)
	}
}
