package compiler_test

import (
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/file"
)

func findOpcodePC(t *testing.T, program *bytecode.Program, opcode bytecode.Opcode) int {
	t.Helper()

	for pc, inst := range program.Bytecode {
		if inst.Opcode == opcode {
			return pc
		}
	}

	t.Fatalf("opcode %s not found", opcode)
	return -1
}

func TestCompiler_OptionalQueryCatchEndsBeforeFollowingInstruction(t *testing.T) {
	source := "LET q = (QUERY ONE `.items` IN @empty USING css)?\nRETURN q.foo"

	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		t.Run(fmt.Sprintf("O%d", level), func(t *testing.T) {
			program, err := compiler.New(compiler.WithOptimizationLevel(level)).Compile(file.NewSource("catch-query", source))
			if err != nil {
				t.Fatalf("compile failed: %v", err)
			}

			if got, want := len(program.CatchTable), 1; got != want {
				t.Fatalf("unexpected catch table size: got %d, want %d", got, want)
			}

			catch := program.CatchTable[0]
			propPC := findOpcodePC(t, program, bytecode.OpLoadPropertyConst)

			if got, want := catch[1], propPC-1; got != want {
				t.Fatalf("unexpected catch end: got %d, want %d", got, want)
			}

			if got, want := catch[2], -1; got != want {
				t.Fatalf("unexpected catch jump: got %d, want %d", got, want)
			}
		})
	}
}

func TestCompiler_OptionalForCatchUsesInclusiveEndAndCleanupJump(t *testing.T) {
	source := "LET xs = (FOR i IN ERROR() RETURN i)?\nRETURN xs.foo"

	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		t.Run(fmt.Sprintf("O%d", level), func(t *testing.T) {
			program, err := compiler.New(compiler.WithOptimizationLevel(level)).Compile(file.NewSource("catch-for", source))
			if err != nil {
				t.Fatalf("compile failed: %v", err)
			}

			if got, want := len(program.CatchTable), 1; got != want {
				t.Fatalf("unexpected catch table size: got %d, want %d", got, want)
			}

			catch := program.CatchTable[0]
			propPC := findOpcodePC(t, program, bytecode.OpLoadPropertyConst)
			closePC := findOpcodePC(t, program, bytecode.OpClose)

			if got, want := catch[1], propPC-1; got != want {
				t.Fatalf("unexpected catch end: got %d, want %d", got, want)
			}

			if got, want := catch[2], closePC; got != want {
				t.Fatalf("unexpected catch jump: got %d, want %d", got, want)
			}
		})
	}
}
