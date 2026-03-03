package optimization_test_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func TestOpExists(t *testing.T) {
	tests := []struct {
		name     string
		program  *bytecode.Program
		expected runtime.Boolean
	}{
		{
			name: "none is false",
			program: programWithOps(
				bytecode.NewInstruction(bytecode.OpLoadNone, bytecode.NewRegister(1)),
				bytecode.NewInstruction(bytecode.OpExists, bytecode.NewRegister(2), bytecode.NewRegister(1)),
				bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(2)),
			),
			expected: runtime.False,
		},
		{
			name: "empty string is false",
			program: programWithConst(runtime.NewString(""),
				bytecode.NewInstruction(bytecode.OpExists, bytecode.NewRegister(2), bytecode.NewRegister(1)),
				bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(2)),
			),
			expected: runtime.False,
		},
		{
			name: "non-empty string is true",
			program: programWithConst(runtime.NewString("ok"),
				bytecode.NewInstruction(bytecode.OpExists, bytecode.NewRegister(2), bytecode.NewRegister(1)),
				bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(2)),
			),
			expected: runtime.True,
		},
		{
			name: "empty array is false",
			program: programWithConst(runtime.EmptyArray(),
				bytecode.NewInstruction(bytecode.OpExists, bytecode.NewRegister(2), bytecode.NewRegister(1)),
				bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(2)),
			),
			expected: runtime.False,
		},
		{
			name: "non-empty array is true",
			program: programWithConst(runtime.NewArrayWith(runtime.NewInt(1)),
				bytecode.NewInstruction(bytecode.OpExists, bytecode.NewRegister(2), bytecode.NewRegister(1)),
				bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(2)),
			),
			expected: runtime.True,
		},
		{
			name: "empty object is false",
			program: programWithConst(runtime.NewObject(),
				bytecode.NewInstruction(bytecode.OpExists, bytecode.NewRegister(2), bytecode.NewRegister(1)),
				bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(2)),
			),
			expected: runtime.False,
		},
		{
			name: "non-empty object is true",
			program: programWithConst(runtime.NewObjectWith(map[string]runtime.Value{"key": runtime.NewInt(1)}),
				bytecode.NewInstruction(bytecode.OpExists, bytecode.NewRegister(2), bytecode.NewRegister(1)),
				bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(2)),
			),
			expected: runtime.True,
		},
		{
			name: "non-measurable value is true",
			program: programWithConst(runtime.NewInt(42),
				bytecode.NewInstruction(bytecode.OpExists, bytecode.NewRegister(2), bytecode.NewRegister(1)),
				bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(2)),
			),
			expected: runtime.True,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			out, err := vm.New(test.program).Run(context.Background(), vm.NewDefaultEnvironment())
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			val, ok := out.(runtime.Boolean)
			if !ok {
				t.Fatalf("expected runtime.Boolean, got %T", out)
			}

			if val != test.expected {
				t.Fatalf("expected %v, got %v", test.expected, val)
			}
		})
	}
}

func programWithConst(value runtime.Value, ops ...bytecode.Instruction) *bytecode.Program {
	instructions := make([]bytecode.Instruction, 0, len(ops)+1)
	instructions = append(instructions, bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(0)))
	instructions = append(instructions, ops...)

	return &bytecode.Program{
		ISAVersion: bytecode.Version,
		Bytecode:   instructions,
		Constants:  []runtime.Value{value},
		Registers:  3,
	}
}

func programWithOps(ops ...bytecode.Instruction) *bytecode.Program {
	return &bytecode.Program{
		ISAVersion: bytecode.Version,
		Bytecode:   ops,
		Registers:  3,
	}
}
