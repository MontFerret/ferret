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
		ProgramCheck(`RETURN "a" + 1 + "b" + 2 + @x + "c" + 3`, func(program *bytecode.Program) error {
			if err := assertOpcodeCount(program.Bytecode, bytecode.OpAdd, 0); err != nil {
				return err
			}

			if err := assertOpcodeCount(program.Bytecode, bytecode.OpConcat, 1); err != nil {
				return err
			}

			return assertProgramStringConstants(program, "a1b2", "c3")
		}, "concat chain merges dynamic literal runs"),
	}, compiler.O0)
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
	}, compiler.O0)
}

func assertOpcodeCount(instructions []bytecode.Instruction, opcode bytecode.Opcode, want int) error {
	got := 0
	for _, inst := range instructions {
		if inst.Opcode == opcode {
			got++
		}
	}

	if got != want {
		return fmt.Errorf("unexpected %s count: got %d, want %d", opcode, got, want)
	}

	return nil
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
