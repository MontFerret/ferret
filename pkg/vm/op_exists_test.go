package vm_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func TestOpExists(t *testing.T) {
	tests := []struct {
		name     string
		program  *vm.Program
		expected runtime.Boolean
	}{
		{
			name: "none is false",
			program: programWithOps(
				vm.NewInstruction(vm.OpLoadNone, vm.NewRegister(1)),
				vm.NewInstruction(vm.OpExists, vm.NewRegister(2), vm.NewRegister(1)),
				vm.NewInstruction(vm.OpReturn, vm.NewRegister(2)),
			),
			expected: runtime.False,
		},
		{
			name: "empty string is false",
			program: programWithConst(runtime.NewString(""),
				vm.NewInstruction(vm.OpExists, vm.NewRegister(2), vm.NewRegister(1)),
				vm.NewInstruction(vm.OpReturn, vm.NewRegister(2)),
			),
			expected: runtime.False,
		},
		{
			name: "non-empty string is true",
			program: programWithConst(runtime.NewString("ok"),
				vm.NewInstruction(vm.OpExists, vm.NewRegister(2), vm.NewRegister(1)),
				vm.NewInstruction(vm.OpReturn, vm.NewRegister(2)),
			),
			expected: runtime.True,
		},
		{
			name: "empty array is false",
			program: programWithConst(runtime.EmptyArray(),
				vm.NewInstruction(vm.OpExists, vm.NewRegister(2), vm.NewRegister(1)),
				vm.NewInstruction(vm.OpReturn, vm.NewRegister(2)),
			),
			expected: runtime.False,
		},
		{
			name: "non-empty array is true",
			program: programWithConst(runtime.NewArrayWith(runtime.NewInt(1)),
				vm.NewInstruction(vm.OpExists, vm.NewRegister(2), vm.NewRegister(1)),
				vm.NewInstruction(vm.OpReturn, vm.NewRegister(2)),
			),
			expected: runtime.True,
		},
		{
			name: "empty object is false",
			program: programWithConst(runtime.NewObject(),
				vm.NewInstruction(vm.OpExists, vm.NewRegister(2), vm.NewRegister(1)),
				vm.NewInstruction(vm.OpReturn, vm.NewRegister(2)),
			),
			expected: runtime.False,
		},
		{
			name: "non-empty object is true",
			program: programWithConst(runtime.NewObjectWith(runtime.NewObjectProperty("a", runtime.NewInt(1))),
				vm.NewInstruction(vm.OpExists, vm.NewRegister(2), vm.NewRegister(1)),
				vm.NewInstruction(vm.OpReturn, vm.NewRegister(2)),
			),
			expected: runtime.True,
		},
		{
			name: "non-measurable value is true",
			program: programWithConst(runtime.NewInt(42),
				vm.NewInstruction(vm.OpExists, vm.NewRegister(2), vm.NewRegister(1)),
				vm.NewInstruction(vm.OpReturn, vm.NewRegister(2)),
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

func programWithConst(value runtime.Value, ops ...vm.Instruction) *vm.Program {
	bytecode := make([]vm.Instruction, 0, len(ops)+1)
	bytecode = append(bytecode, vm.NewInstruction(vm.OpLoadConst, vm.NewRegister(1), vm.NewConstant(0)))
	bytecode = append(bytecode, ops...)
	return &vm.Program{
		Bytecode:  bytecode,
		Constants: []runtime.Value{value},
		Registers: 3,
	}
}

func programWithOps(ops ...vm.Instruction) *vm.Program {
	return &vm.Program{
		Bytecode:  ops,
		Registers: 3,
	}
}
