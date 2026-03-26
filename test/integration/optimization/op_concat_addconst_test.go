package optimization_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/test/spec"
	"github.com/MontFerret/ferret/v2/test/spec/assert"
)

func TestOpConcat(t *testing.T) {
	testCases := []struct {
		name     string
		expected string
		values   []runtime.Value
		startReg int
		count    int
	}{
		{
			name:     "count zero returns empty string",
			values:   nil,
			startReg: 1,
			count:    0,
			expected: "",
		},
		{
			name: "count one coerces to string",
			values: []runtime.Value{
				runtime.NewInt(42),
			},
			startReg: 1,
			count:    1,
			expected: "42",
		},
		{
			name: "count two fast-path skips empty prefix",
			values: []runtime.Value{
				runtime.EmptyString,
				runtime.NewString("bar"),
			},
			startReg: 1,
			count:    2,
			expected: "bar",
		},
		{
			name: "count three fast-path handles trailing empties",
			values: []runtime.Value{
				runtime.NewString("a"),
				runtime.EmptyString,
				runtime.EmptyString,
			},
			startReg: 1,
			count:    3,
			expected: "a",
		},
		{
			name: "builder path with all empty parts returns empty string",
			values: []runtime.Value{
				runtime.None,
				runtime.EmptyString,
				runtime.None,
				runtime.EmptyString,
			},
			startReg: 1,
			count:    4,
			expected: "",
		},
		{
			name: "builder path with unicode and mixed scalar types",
			values: []runtime.Value{
				runtime.NewString("пр"),
				runtime.NewInt(7),
				runtime.True,
				runtime.NewString("!"),
			},
			startReg: 5,
			count:    4,
			expected: "пр7true!",
		},
	}

	specs := make([]spec.Spec, 0, len(testCases))
	for _, tc := range testCases {
		specs = append(specs, spec.NewWith(spec.NewProgramInput(buildConcatProgram(tc.values, tc.startReg, tc.count)), tc.name).Expect().Exec(assert.ShouldEqual, tc.expected))
	}

	runProgramSpecs(t, specs)
}

func TestOpAddConst(t *testing.T) {
	testCases := []struct {
		left     runtime.Value
		right    runtime.Value
		expected any
		name     string
	}{
		{
			name:     "int plus int",
			left:     runtime.NewInt(1),
			right:    runtime.NewInt(2),
			expected: 3,
		},
		{
			name:     "float plus int",
			left:     runtime.NewFloat(1.5),
			right:    runtime.NewInt(2),
			expected: 3.5,
		},
		{
			name:     "int plus string",
			left:     runtime.NewInt(1),
			right:    runtime.NewString("x"),
			expected: "1x",
		},
		{
			name:     "string plus int",
			left:     runtime.NewString("x"),
			right:    runtime.NewInt(1),
			expected: "x1",
		},
	}

	specs := make([]spec.Spec, 0, len(testCases))
	for _, tc := range testCases {
		program := &bytecode.Program{
			ISAVersion: bytecode.Version,
			Constants: []runtime.Value{
				tc.left,
				tc.right,
			},
			Bytecode: []bytecode.Instruction{
				bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(0)),
				bytecode.NewInstruction(bytecode.OpAddConst, bytecode.NewRegister(2), bytecode.NewRegister(1), bytecode.NewConstant(1)),
				bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(2)),
			},
			Registers: 3,
		}

		specs = append(specs, spec.NewWith(spec.NewProgramInput(program), tc.name).Expect().Exec(assert.ShouldEqual, tc.expected))
	}

	runProgramSpecs(t, specs)
}

func buildConcatProgram(values []runtime.Value, startReg, count int) *bytecode.Program {
	instructions := make([]bytecode.Instruction, 0, len(values)+2)

	for i := 0; i < len(values); i++ {
		instructions = append(instructions, bytecode.NewInstruction(
			bytecode.OpLoadConst,
			bytecode.NewRegister(startReg+i),
			bytecode.NewConstant(i),
		))
	}

	dst := startReg + len(values) + 1

	instructions = append(instructions, bytecode.NewInstruction(
		bytecode.OpConcat,
		bytecode.NewRegister(dst),
		bytecode.NewRegister(startReg),
		bytecode.Operand(count),
	))
	instructions = append(instructions, bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(dst)))

	return &bytecode.Program{
		ISAVersion: bytecode.Version,
		Constants:  values,
		Bytecode:   instructions,
		Registers:  dst + 1,
	}
}
