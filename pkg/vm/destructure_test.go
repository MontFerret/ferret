package vm

import (
	"context"
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestAssertDestructureExecution(t *testing.T) {
	tests := []struct {
		value runtime.Value
		name  string
		mode  bytecode.DestructureMode
	}{
		{
			name:  "object",
			mode:  bytecode.DestructureModeObject,
			value: runtime.NewObject(),
		},
		{
			name:  "array",
			mode:  bytecode.DestructureModeArray,
			value: runtime.NewArray(0),
		},
		{
			name:  "none as object",
			mode:  bytecode.DestructureModeObject,
			value: runtime.None,
		},
		{
			name:  "none as array",
			mode:  bytecode.DestructureModeArray,
			value: runtime.None,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			program := assertDestructureProgram(test.value, test.mode)
			instance := mustNewVM(t, program)
			result, err := instance.Run(context.Background(), mustNewEnvironment(t))
			if err != nil {
				t.Fatal(err)
			}

			if got := result.Root(); got != test.value {
				t.Fatalf("result = %v, want original value", got)
			}

			if err := result.Close(); err != nil {
				t.Fatalf("close result: %v", err)
			}
		})
	}
}

func TestAssertDestructureExecutionRejectsWrongShape(t *testing.T) {
	program := assertDestructureProgram(runtime.NewInt(1), bytecode.DestructureModeObject)
	instance := mustNewVM(t, program)

	_, err := instance.Run(context.Background(), mustNewEnvironment(t))
	if err == nil {
		t.Fatal("expected destructuring error")
	}

	var runtimeErr *RuntimeError
	if !errors.As(err, &runtimeErr) {
		t.Fatalf("error = %T, want *RuntimeError", err)
	}

	if got, want := runtimeErr.Message, "cannot destructure Int as Object"; got != want {
		t.Fatalf("message = %q, want %q", got, want)
	}
}

func TestRegexpWarmupPreservesRegisterFactsAcrossDestructureAssertion(t *testing.T) {
	const regexpPC = 2

	program := newTestProgram(
		3,
		[]runtime.Value{runtime.NewString("value")},
		bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(2), bytecode.NewConstant(0)),
		bytecode.NewInstruction(bytecode.OpAssertDestructure, bytecode.NewRegister(2), bytecode.Operand(bytecode.DestructureModeArray)),
		bytecode.NewInstruction(bytecode.OpRegexp, bytecode.NewRegister(1), bytecode.NewRegister(0), bytecode.NewRegister(2)),
		bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
	)
	instance := mustNewVM(t, program)

	if err := ensureRegexpsWarmed(instance); err != nil {
		t.Fatal(err)
	}

	if instance.cache.Regexps[regexpPC] == nil {
		t.Fatal("expected regexp after destructure assertion to be warmed")
	}
}

func assertDestructureProgram(value runtime.Value, mode bytecode.DestructureMode) *bytecode.Program {
	return &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  1,
		Constants:  []runtime.Value{value},
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(0), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpAssertDestructure, bytecode.NewRegister(0), bytecode.Operand(mode)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
		},
	}
}
