package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
)

func findFirstOpcodeIndex(code []bytecode.Instruction, op bytecode.Opcode) (int, bool) {
	for i, inst := range code {
		if inst.Opcode == op {
			return i, true
		}
	}

	return -1, false
}

func hasOpcode(code []bytecode.Instruction, op bytecode.Opcode) bool {
	_, ok := findFirstOpcodeIndex(code, op)
	return ok
}

func lastRegisterDefOpcodeBefore(code []bytecode.Instruction, before int, reg int) (bytecode.Opcode, bool) {
	for i := before - 1; i >= 0; i-- {
		inst := code[i]
		if !inst.Operands[0].IsRegister() {
			continue
		}
		if inst.Operands[0].Register() == reg {
			return inst.Opcode, true
		}
	}

	return bytecode.OpMove, false
}

func assertCallArgsLoadedFromConsts(t *testing.T, code []bytecode.Instruction, callIndex int, expectedArgs int) {
	t.Helper()

	call := code[callIndex]
	if !call.Operands[1].IsRegister() || !call.Operands[2].IsRegister() {
		t.Fatalf("expected register range operands in call, got %v", call.Operands)
	}

	start := call.Operands[1].Register()
	end := call.Operands[2].Register()
	if got := end - start + 1; got != expectedArgs {
		t.Fatalf("expected %d call args, got %d (range R%d..R%d)", expectedArgs, got, start, end)
	}

	for reg := start; reg <= end; reg++ {
		op, ok := lastRegisterDefOpcodeBefore(code, callIndex, reg)
		if !ok {
			t.Fatalf("expected to find definition for argument register R%d", reg)
		}

		if op != bytecode.OpLoadConst {
			t.Fatalf("expected argument register R%d to be loaded via LOADC, got %s", reg, op)
		}
	}
}

func TestUdfCallConstantArgsDirectLoadO0(t *testing.T) {
	expr := `
FUNC f2(x, y) => x + y
RETURN f2(1, 2)
`

	prog := compileWithLevel(t, compiler.O0, expr)

	callIndex, ok := findFirstOpcodeIndex(prog.Bytecode, bytecode.OpCall)
	if !ok {
		t.Fatalf("expected OpCall in bytecode")
	}

	assertCallArgsLoadedFromConsts(t, prog.Bytecode, callIndex, 2)

	if hasOpcode(prog.Bytecode, bytecode.OpMove) || hasOpcode(prog.Bytecode, bytecode.OpMoveTracked) {
		t.Fatalf("expected no MOVE/MOVET instructions for constant-only UDF call setup")
	}
}

func TestHostCallConstantArgsDirectLoadO0(t *testing.T) {
	expr := `RETURN TEST(1, 2)`

	prog := compileWithLevel(t, compiler.O0, expr)

	callIndex, ok := findFirstOpcodeIndex(prog.Bytecode, bytecode.OpHCall)
	if !ok {
		t.Fatalf("expected OpHCall in bytecode")
	}

	assertCallArgsLoadedFromConsts(t, prog.Bytecode, callIndex, 2)

	if hasOpcode(prog.Bytecode, bytecode.OpMove) || hasOpcode(prog.Bytecode, bytecode.OpMoveTracked) {
		t.Fatalf("expected no MOVE/MOVET instructions for constant-only host call setup")
	}
}

func TestCallArgumentLoweringKeepsMoveForNonLiteralArgO0(t *testing.T) {
	expr := `
LET x = 1
RETURN TEST(x, 2)
`

	prog := compileWithLevel(t, compiler.O0, expr)

	if !hasOpcode(prog.Bytecode, bytecode.OpMoveTracked) {
		t.Fatalf("expected MOVET instruction for non-literal argument setup")
	}
}
