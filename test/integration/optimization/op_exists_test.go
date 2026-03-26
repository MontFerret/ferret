package optimization_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/test/spec"
	"github.com/MontFerret/ferret/v2/test/spec/assert"
)

func TestOpExists(t *testing.T) {
	tests := []struct {
		program  *bytecode.Program
		name     string
		expected bool
	}{
		{
			name: "none is false",
			program: programWithOps(
				bytecode.NewInstruction(bytecode.OpLoadNone, bytecode.NewRegister(1)),
				bytecode.NewInstruction(bytecode.OpExists, bytecode.NewRegister(2), bytecode.NewRegister(1)),
				bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(2)),
			),
			expected: false,
		},
		{
			name: "empty string is false",
			program: programWithConst(runtime.NewString(""),
				bytecode.NewInstruction(bytecode.OpExists, bytecode.NewRegister(2), bytecode.NewRegister(1)),
				bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(2)),
			),
			expected: false,
		},
		{
			name: "non-empty string is true",
			program: programWithConst(runtime.NewString("ok"),
				bytecode.NewInstruction(bytecode.OpExists, bytecode.NewRegister(2), bytecode.NewRegister(1)),
				bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(2)),
			),
			expected: true,
		},
		{
			name: "empty array is false",
			program: programWithConst(runtime.EmptyArray(),
				bytecode.NewInstruction(bytecode.OpExists, bytecode.NewRegister(2), bytecode.NewRegister(1)),
				bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(2)),
			),
			expected: false,
		},
		{
			name: "non-empty array is true",
			program: programWithConst(runtime.NewArrayWith(runtime.NewInt(1)),
				bytecode.NewInstruction(bytecode.OpExists, bytecode.NewRegister(2), bytecode.NewRegister(1)),
				bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(2)),
			),
			expected: true,
		},
		{
			name: "empty object is false",
			program: programWithConst(runtime.NewObject(),
				bytecode.NewInstruction(bytecode.OpExists, bytecode.NewRegister(2), bytecode.NewRegister(1)),
				bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(2)),
			),
			expected: false,
		},
		{
			name: "non-empty object is true",
			program: programWithConst(runtime.NewObjectWith(map[string]runtime.Value{"key": runtime.NewInt(1)}),
				bytecode.NewInstruction(bytecode.OpExists, bytecode.NewRegister(2), bytecode.NewRegister(1)),
				bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(2)),
			),
			expected: true,
		},
		{
			name: "non-measurable value is true",
			program: programWithConst(runtime.NewInt(42),
				bytecode.NewInstruction(bytecode.OpExists, bytecode.NewRegister(2), bytecode.NewRegister(1)),
				bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(2)),
			),
			expected: true,
		},
	}

	specs := make([]spec.Spec, 0, len(tests))
	for _, test := range tests {
		specs = append(specs, spec.NewWith(spec.NewProgramInput(test.program), test.name).Expect().Exec(assert.ShouldEqual, test.expected))
	}

	runProgramSpecs(t, specs)
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
