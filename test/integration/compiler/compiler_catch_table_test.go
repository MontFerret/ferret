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

func expectOpcodeAbsent(program *bytecode.Program, opcode bytecode.Opcode) error {
	if pc, err := findOpcodePC(program, opcode); err == nil {
		return fmt.Errorf("unexpected opcode %s at pc %d", opcode, pc)
	}

	return nil
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

func TestCompiler_RetryWithoutFallbackLeavesFinalAttemptUnprotected(t *testing.T) {
	RunSpecsLevels(t, []spec.Spec{
		ProgramCheck("RETURN FAIL() ON ERROR RETRY 2", func(program *bytecode.Program) error {
			if got, want := len(program.CatchTable), 1; got != want {
				return fmt.Errorf("unexpected catch table size: got %d, want %d", got, want)
			}

			callCount := 0
			for _, inst := range program.Bytecode {
				if inst.Opcode == bytecode.OpHCall {
					callCount++
				}
			}

			if got, want := callCount, 2; got != want {
				return fmt.Errorf("unexpected host call count: got %d, want %d", got, want)
			}

			return nil
		}, "retry without fallback should leave one final uncaught attempt"),
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

func TestCompiler_GroupedForRetryRoutesThroughCleanup(t *testing.T) {
	RunSpecsLevels(t, []spec.Spec{
		ProgramCheck("LET xs = (FOR i IN [1, 2] LET y = ERROR() RETURN y + i) ON ERROR RETRY 2 OR RETURN []\nRETURN LENGTH(xs)", func(program *bytecode.Program) error {
			if got, want := len(program.CatchTable), 1; got != want {
				return fmt.Errorf("unexpected catch table size: got %d, want %d", got, want)
			}

			catch := program.CatchTable[0]
			closePC, err := findOpcodePC(program, bytecode.OpClose)
			if err != nil {
				return err
			}

			lengthPC, err := findOpcodePC(program, bytecode.OpLength)
			if err != nil {
				return err
			}

			if got := catch[2]; got <= closePC || got >= lengthPC {
				return fmt.Errorf("unexpected grouped for retry handler pc: got %d, want (%d, %d)", got, closePC, lengthPC)
			}

			return nil
		}, "grouped for retry should route failures through loop cleanup"),
	}, compiler.O0, compiler.O1)
}

func TestCompiler_GroupedForRecoveryRoutesThroughCleanup(t *testing.T) {
	RunSpecsLevels(t, []spec.Spec{
		ProgramCheck("LET xs = (FOR i IN [1, 2] LET y = ERROR() RETURN y + i) ON ERROR RETURN []\nRETURN LENGTH(xs)", func(program *bytecode.Program) error {
			if got, want := len(program.CatchTable), 1; got != want {
				return fmt.Errorf("unexpected catch table size: got %d, want %d", got, want)
			}

			catch := program.CatchTable[0]
			closePC, err := findOpcodePC(program, bytecode.OpClose)
			if err != nil {
				return err
			}

			lengthPC, err := findOpcodePC(program, bytecode.OpLength)
			if err != nil {
				return err
			}

			if got := catch[2]; got <= closePC || got >= lengthPC {
				return fmt.Errorf("unexpected grouped for recovery handler pc: got %d, want (%d, %d)", got, closePC, lengthPC)
			}

			jumpPC := catch[2] + 1
			if jumpPC >= len(program.Bytecode) {
				return fmt.Errorf("grouped for recovery jump pc %d out of range", jumpPC)
			}

			if got, want := program.Bytecode[jumpPC].Opcode, bytecode.OpJump; got != want {
				return fmt.Errorf("unexpected grouped for recovery opcode at pc %d: got %s, want %s", jumpPC, got, want)
			}

			if got, want := int(program.Bytecode[jumpPC].Operands[0]), closePC; got != want {
				return fmt.Errorf("unexpected grouped for recovery cleanup jump: got %d, want %d", got, want)
			}

			return nil
		}, "grouped for recovery should route failures through loop cleanup"),
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

func TestCompiler_WaitForEventRetryUsesCleanupJump(t *testing.T) {
	RunSpecsLevels(t, []spec.Spec{
		ProgramCheck("LET ok = WAITFOR EVENT \"test\" IN @obs ON ERROR RETRY 2 OR RETURN NONE\nRETURN ok.foo", func(program *bytecode.Program) error {
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
				return fmt.Errorf("unexpected waitfor event retry jump: got %d, want (%d, %d)", got, closePC, propPC)
			}

			return nil
		}, "waitfor event retry should route failures through stream cleanup"),
	}, compiler.O0, compiler.O1)
}

func TestCompiler_GroupedWaitForEventRecoveryUsesTimeoutAwarePath(t *testing.T) {
	RunSpecsLevels(t, []spec.Spec{
		ProgramCheck("RETURN (WAITFOR EVENT \"test\" IN @obs TIMEOUT 1ms) ON TIMEOUT RETURN \"timeout\" ON ERROR RETURN \"error\"", func(program *bytecode.Program) error {
			if _, err := findOpcodePC(program, bytecode.OpIterNextTimeout); err != nil {
				return err
			}

			return expectOpcodeAbsent(program, bytecode.OpIterNext)
		}, "grouped waitfor event recovery should use timeout-aware bytecode"),
		ProgramCheck("RETURN (WAITFOR EVENT \"test\" IN @obs TIMEOUT 1ms) ON TIMEOUT RETURN \"timeout\" ON ERROR RETRY 2 OR RETURN \"error\"", func(program *bytecode.Program) error {
			if _, err := findOpcodePC(program, bytecode.OpIterNextTimeout); err != nil {
				return err
			}

			return expectOpcodeAbsent(program, bytecode.OpIterNext)
		}, "grouped waitfor event retry should use timeout-aware bytecode"),
	}, compiler.O0, compiler.O1)
}

func TestCompiler_RecoveryReturnWideningKeepsGenericMemberAccess(t *testing.T) {
	RunSpecsLevels(t, []spec.Spec{
		ProgramCheck("LET x = QUERY VALUE `.items` IN @doc USING css ON ERROR RETURN { \"0\": \"fallback\" }\nRETURN x[0]", func(program *bytecode.Program) error {
			lastLoadPC := -1
			lastLoadOpcode := bytecode.Opcode(0)

			for pc, inst := range program.Bytecode {
				switch inst.Opcode {
				case bytecode.OpLoadPropertyConst, bytecode.OpLoadKeyConst, bytecode.OpLoadIndexConst:
					lastLoadPC = pc
					lastLoadOpcode = inst.Opcode
				}
			}

			if lastLoadPC < 0 {
				return fmt.Errorf("no member load opcode found")
			}

			if got, want := lastLoadOpcode, bytecode.OpLoadPropertyConst; got != want {
				return fmt.Errorf("unexpected final member load opcode at pc %d: got %s, want %s", lastLoadPC, got, want)
			}

			return nil
		}, "recovery return should widen result type for later member access"),
	}, compiler.O0, compiler.O1)
}
