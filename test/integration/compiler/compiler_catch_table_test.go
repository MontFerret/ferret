package compiler_test

import (
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/compile"
)

func findOpcodePC(program *bytecode.Program, opcode bytecode.Opcode) (int, error) {
	for pc, inst := range program.Bytecode {
		if inst.Opcode == opcode {
			return pc, nil
		}
	}

	return -1, fmt.Errorf("opcode %s not found", opcode)
}

func TestCompiler_OptionalQueryCatchEndsBeforeFollowingInstruction(t *testing.T) {
	RunSpecsLevels(t, []spec.Spec{
		ProgramCheck("LET q = (QUERY ONE `.items` IN @empty USING css)?\nRETURN q.foo", func(program *bytecode.Program) error {
			if got, want := len(program.CatchTable), 1; got != want {
				return fmt.Errorf("unexpected catch table size: got %d, want %d", got, want)
			}

			catch := program.CatchTable[0]
			propPC, err := findOpcodePC(program, bytecode.OpLoadPropertyConst)
			if err != nil {
				return err
			}

			if got, want := catch[1], propPC-3; got != want {
				return fmt.Errorf("unexpected catch end: got %d, want %d", got, want)
			}

			if got, want := catch[2], propPC-1; got != want {
				return fmt.Errorf("unexpected catch jump: got %d, want %d", got, want)
			}

			return nil
		}, "optional query catch ends before following instruction"),
	}, compiler.O0, compiler.O1)
}

func TestCompiler_ExplicitQuerySuppressCatchEndsBeforeFollowingInstruction(t *testing.T) {
	RunSpecsLevels(t, []spec.Spec{
		ProgramCheck("LET q = QUERY ONE `.items` IN @empty USING css ON ERROR RETURN NONE\nRETURN q.foo", func(program *bytecode.Program) error {
			if got, want := len(program.CatchTable), 1; got != want {
				return fmt.Errorf("unexpected catch table size: got %d, want %d", got, want)
			}

			catch := program.CatchTable[0]
			propPC, err := findOpcodePC(program, bytecode.OpLoadPropertyConst)
			if err != nil {
				return err
			}

			if got, want := catch[1], propPC-3; got != want {
				return fmt.Errorf("unexpected catch end: got %d, want %d", got, want)
			}

			if got, want := catch[2], propPC-1; got != want {
				return fmt.Errorf("unexpected catch jump: got %d, want %d", got, want)
			}

			return nil
		}, "explicit query suppress catch ends before following instruction"),
	}, compiler.O0, compiler.O1)
}

func TestCompiler_OptionalForCatchUsesInclusiveEndAndCleanupJump(t *testing.T) {
	RunSpecsLevels(t, []spec.Spec{
		ProgramCheck("LET xs = (FOR i IN ERROR() RETURN i)?\nRETURN xs.foo", func(program *bytecode.Program) error {
			if got, want := len(program.CatchTable), 1; got != want {
				return fmt.Errorf("unexpected catch table size: got %d, want %d", got, want)
			}

			catch := program.CatchTable[0]
			propPC, err := findOpcodePC(program, bytecode.OpLoadPropertyConst)
			if err != nil {
				return err
			}
			closePC, err := findOpcodePC(program, bytecode.OpClose)
			if err != nil {
				return err
			}

			if got, want := catch[1], propPC-1; got != want {
				return fmt.Errorf("unexpected catch end: got %d, want %d", got, want)
			}

			if got, want := catch[2], closePC; got != want {
				return fmt.Errorf("unexpected catch jump: got %d, want %d", got, want)
			}

			return nil
		}, "optional for catch uses cleanup jump"),
	}, compiler.O0, compiler.O1)
}

func TestCompiler_WaitForEventSuppressCatchUsesCleanupJump(t *testing.T) {
	RunSpecsLevels(t, []spec.Spec{
		ProgramCheck("LET ok = WAITFOR EVENT \"test\" IN @obs ON ERROR RETURN NONE\nRETURN ok.foo", func(program *bytecode.Program) error {
			if got, want := len(program.CatchTable), 1; got != want {
				return fmt.Errorf("unexpected catch table size: got %d, want %d", got, want)
			}

			catch := program.CatchTable[0]
			propPC, err := findOpcodePC(program, bytecode.OpLoadPropertyConst)
			if err != nil {
				return err
			}
			closePC, err := findOpcodePC(program, bytecode.OpClose)
			if err != nil {
				return err
			}

			if got, want := catch[1], closePC; got != want {
				return fmt.Errorf("unexpected catch end: got %d, want %d", got, want)
			}

			if got := catch[2]; got <= closePC || got >= propPC {
				return fmt.Errorf("unexpected waitfor event recovery jump: got %d, want (%d, %d)", got, closePC, propPC)
			}

			return nil
		}, "waitfor event suppress catch uses cleanup jump"),
	}, compiler.O0, compiler.O1)
}
