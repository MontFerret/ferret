package compiler_test

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/compile"
)

func TestCompilerConcatChainMergesDynamicLiteralRuns(t *testing.T) {
	RunSpecsLevels(t, []spec.Spec{
		ProgramCheck(`LET x = "x"
RETURN "a" + 1 + "b" + 2 + x + "c" + 3`, func(program *bytecode.Program) error {
			if err := assertOpcodeCount(program.Bytecode, bytecode.OpAdd, 0); err != nil {
				return err
			}

			if err := assertOpcodeCount(program.Bytecode, bytecode.OpConcat, 1); err != nil {
				return err
			}

			return assertProgramStringConstants(program, "a1b2", "c3", "x")
		}, "concat chain merges dynamic literal runs"),
	}, compiler.OptimizationNone)
}

func TestCompilerStringAnchorsAcceptTemporalAndDynamicValues(t *testing.T) {
	tests := []struct {
		name string
		src  string
	}{
		{"DateTime plus String", `RETURN NOW() + "5m"`},
		{"dynamic value plus String", `RETURN @value + "500ms"`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, level := range []compiler.OptimizationLevel{compiler.OptimizationNone, compiler.OptimizationFull} {
				program := compileWithLevel(t, level, test.src)

				if err := assertOpcodeCount(program.Bytecode, bytecode.OpConcat, 0); err != nil {
					t.Fatal(err)
				}

				stringAdds := countOpcodes(program.Bytecode, bytecode.OpAdd) + countOpcodes(program.Bytecode, bytecode.OpAddConst)
				if stringAdds != 1 {
					t.Fatalf("O%d: unexpected String-add count: got %d, want 1", level, stringAdds)
				}
			}
		})
	}
}

func TestCompilerKnownStringAdditionRetainsConcatOptimization(t *testing.T) {
	for _, level := range []compiler.OptimizationLevel{compiler.OptimizationNone, compiler.OptimizationFull} {
		program := compileWithLevel(t, level, `LET prefix = "a"
RETURN prefix + "b"`)

		if got := countOpcodes(program.Bytecode, bytecode.OpAdd); got != 0 {
			t.Fatalf("O%d: unexpected %s count: got %d, want 0", level, bytecode.OpAdd, got)
		}
		if got := countOpcodes(program.Bytecode, bytecode.OpConcat); got != 0 {
			t.Fatalf("O%d: unexpected %s count: got %d, want 0", level, bytecode.OpConcat, got)
		}
		if level == compiler.OptimizationNone && countOpcodes(program.Bytecode, bytecode.OpAddConst) != 1 {
			t.Fatalf("O%d: known string addition did not use %s", level, bytecode.OpAddConst)
		}
	}
}

func TestCompilerInvalidTemporalOrderingIsNotConstantFolded(t *testing.T) {
	program := compileWithLevel(t, compiler.OptimizationFull, `RETURN 1s < "tomorrow"`)
	if got := countOpcodes(program.Bytecode, bytecode.OpLt); got != 1 {
		t.Fatalf("unexpected %s count: got %d, want 1", bytecode.OpLt, got)
	}
}

func TestCompilerStringConcatAssignmentUsesMergedSegments(t *testing.T) {
	RunSpecsLevels(t, []spec.Spec{
		ProgramCheck(`VAR str = ""

str += " " + 1 + " " + 2 + " " + 3 + " " + 4 + " " + 5

RETURN str`, func(program *bytecode.Program) error {
			if err := assertOpcodeCount(program.Bytecode, bytecode.OpAdd, 0); err != nil {
				return err
			}

			if err := assertOpcodeCount(program.Bytecode, bytecode.OpAddConst, 1); err != nil {
				return err
			}

			if err := assertOpcodeCount(program.Bytecode, bytecode.OpConcat, 0); err != nil {
				return err
			}

			return assertProgramStringConstants(program, "", " 1 2 3 4 5")
		}, "string concat assignment uses merged segments"),
	}, compiler.OptimizationNone)
}

func assertOpcodeCount(instructions []bytecode.Instruction, opcode bytecode.Opcode, want int) error {
	got := countOpcodes(instructions, opcode)

	if got != want {
		return fmt.Errorf("unexpected %s count: got %d, want %d", opcode, got, want)
	}

	return nil
}

func countOpcodes(instructions []bytecode.Instruction, opcode bytecode.Opcode) int {
	count := 0
	for _, inst := range instructions {
		if inst.Opcode == opcode {
			count++
		}
	}

	return count
}

func assertProgramStringConstants(program *bytecode.Program, want ...string) error {
	got := make([]string, 0, len(program.Constants))

	for _, value := range program.Constants {
		str, ok := value.(runtime.String)
		if !ok {
			return fmt.Errorf("unexpected constant type %T in program constants", value)
		}

		got = append(got, str.String())
	}

	sort.Strings(got)
	sort.Strings(want)

	if !reflect.DeepEqual(got, want) {
		return fmt.Errorf("unexpected constants: got %v, want %v", got, want)
	}

	return nil
}
