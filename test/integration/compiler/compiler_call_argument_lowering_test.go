package compiler_test

import (
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/test/spec"
)

func callArgsLoadedFromConsts(code []bytecode.Instruction, callIndex int, expectedArgs int) error {
	call := code[callIndex]
	if !call.Operands[1].IsRegister() || !call.Operands[2].IsRegister() {
		return fmt.Errorf("expected register range operands in call, got %v", call.Operands)
	}

	start := call.Operands[1].Register()
	end := call.Operands[2].Register()
	if got := end - start + 1; got != expectedArgs {
		return fmt.Errorf("expected %d call args, got %d (range R%d..R%d)", expectedArgs, got, start, end)
	}

	for reg := start; reg <= end; reg++ {
		op, ok := lastRegisterDefOpcodeBefore(code, callIndex, reg)
		if !ok {
			return fmt.Errorf("expected to find definition for argument register R%d", reg)
		}

		if op != bytecode.OpLoadConst {
			return fmt.Errorf("expected argument register R%d to be loaded via LOADC, got %s", reg, op)
		}
	}

	return nil
}

func TestUdfCallConstantArgsDirectLoadO0(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		ProgramCheck(`
FUNC f2(x, y) => x + y
RETURN f2(1, 2)
`, func(prog *bytecode.Program) error {
			callIndex, ok := findFirstOpcodeIndex(prog.Bytecode, bytecode.OpCall)
			if !ok {
				return fmt.Errorf("expected OpCall in bytecode")
			}
			if err := callArgsLoadedFromConsts(prog.Bytecode, callIndex, 2); err != nil {
				return err
			}
			if hasOpcode(prog.Bytecode, bytecode.OpMove) || hasOpcode(prog.Bytecode, bytecode.OpMoveTracked) {
				return fmt.Errorf("expected no MOVE/MOVET instructions for constant-only UDF call setup")
			}

			return nil
		}, "udf constant args direct load"),
	})
}

func TestHostCallConstantArgsDirectLoadO0(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		ProgramCheck(`RETURN TEST(1, 2)`, func(prog *bytecode.Program) error {
			callIndex, ok := findFirstOpcodeIndex(prog.Bytecode, bytecode.OpHCall)
			if !ok {
				return fmt.Errorf("expected OpHCall in bytecode")
			}
			if err := callArgsLoadedFromConsts(prog.Bytecode, callIndex, 2); err != nil {
				return err
			}
			if hasOpcode(prog.Bytecode, bytecode.OpMove) || hasOpcode(prog.Bytecode, bytecode.OpMoveTracked) {
				return fmt.Errorf("expected no MOVE/MOVET instructions for constant-only host call setup")
			}

			return nil
		}, "host constant args direct load"),
	})
}

func TestCallArgumentLoweringKeepsMoveForNonLiteralArgO0(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		ProgramCheck(`
LET x = 1
RETURN TEST(x, 2)
`, func(prog *bytecode.Program) error {
			if !hasOpcode(prog.Bytecode, bytecode.OpMoveTracked) {
				return fmt.Errorf("expected MOVET instruction for non-literal argument setup")
			}

			return nil
		}, "non-literal arg keeps tracked move"),
	})
}
