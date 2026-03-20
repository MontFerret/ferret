package compiler

import (
	"reflect"
	"sort"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/file"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestCompilerConcatChainMergesDynamicLiteralRuns(t *testing.T) {
	c := New(WithOptimizationLevel(O0))

	program, err := c.Compile(file.NewAnonymousSource(`RETURN "a" + 1 + "b" + 2 + @x + "c" + 3`))
	if err != nil {
		t.Fatal(err)
	}

	assertOpcodeCount(t, program.Bytecode, bytecode.OpAdd, 0)
	assertOpcodeCount(t, program.Bytecode, bytecode.OpConcat, 1)
	assertProgramStringConstants(t, program, "a1b2", "c3")
}

func TestCompilerStringConcatAssignmentUsesMergedSegments(t *testing.T) {
	c := New(WithOptimizationLevel(O0))

	program, err := c.Compile(file.NewAnonymousSource(`VAR str = ""

str += " " + 1 + " " + 2 + " " + 3 + " " + 4 + " " + 5

RETURN str`))
	if err != nil {
		t.Fatal(err)
	}

	assertOpcodeCount(t, program.Bytecode, bytecode.OpAdd, 0)
	assertOpcodeCount(t, program.Bytecode, bytecode.OpAddConst, 1)
	assertOpcodeCount(t, program.Bytecode, bytecode.OpConcat, 0)
	assertProgramStringConstants(t, program, "", " 1 2 3 4 5")
}

func assertOpcodeCount(t *testing.T, instructions []bytecode.Instruction, opcode bytecode.Opcode, want int) {
	t.Helper()

	got := 0
	for _, inst := range instructions {
		if inst.Opcode == opcode {
			got++
		}
	}

	if got != want {
		t.Fatalf("unexpected %s count: got %d, want %d", opcode, got, want)
	}
}

func assertProgramStringConstants(t *testing.T, program *bytecode.Program, want ...string) {
	t.Helper()

	got := make([]string, 0, len(program.Constants))

	for _, value := range program.Constants {
		str, ok := value.(runtime.String)
		if !ok {
			t.Fatalf("unexpected constant type %T in program constants", value)
		}

		got = append(got, str.String())
	}

	sort.Strings(got)
	sort.Strings(want)

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected constants: got %v, want %v", got, want)
	}
}
