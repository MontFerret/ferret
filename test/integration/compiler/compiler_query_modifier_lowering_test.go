package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func findApplyQueryDescriptorSize(code []bytecode.Instruction, applyIdx int) (int, bool) {
	if applyIdx < 0 || applyIdx >= len(code) {
		return 0, false
	}

	queryReg := code[applyIdx].Operands[2]
	if !queryReg.IsRegister() {
		return 0, false
	}

	for i := applyIdx - 1; i >= 0; i-- {
		inst := code[i]
		if inst.Opcode != bytecode.OpLoadArray {
			continue
		}

		if !inst.Operands[0].IsRegister() || inst.Operands[0].Register() != queryReg.Register() {
			continue
		}

		return int(inst.Operands[1]), true
	}

	return 0, false
}

func assertThreeSlotQueryDescriptor(t *testing.T, code []bytecode.Instruction) {
	t.Helper()

	applyIdx, ok := findFirstOpcodeIndex(code, bytecode.OpApplyQuery)
	if !ok {
		t.Fatalf("expected OpApplyQuery in bytecode")
	}

	size, ok := findApplyQueryDescriptorSize(code, applyIdx)
	if !ok {
		t.Fatalf("expected OpLoadArray for query descriptor before OpApplyQuery")
	}

	if size != 3 {
		t.Fatalf("expected 3-slot query descriptor, got %d", size)
	}
}

func assertFailPrelude(t *testing.T, prog *bytecode.Program, expectedMessage runtime.String) {
	t.Helper()

	failIdx, ok := findFirstOpcodeIndex(prog.Bytecode, bytecode.OpFail)
	if !ok {
		t.Fatalf("expected OpFail in bytecode")
	}

	if failIdx == 0 {
		t.Fatalf("expected OpLoadNone before OpFail")
	}

	if got := prog.Bytecode[failIdx-1].Opcode; got != bytecode.OpLoadNone {
		t.Fatalf("expected OpLoadNone before OpFail, got %s", got)
	}

	fail := prog.Bytecode[failIdx]
	if !fail.Operands[0].IsConstant() {
		t.Fatalf("expected OpFail to use constant-string payload")
	}

	msgIdx := fail.Operands[0].Constant()
	if msgIdx < 0 || msgIdx >= len(prog.Constants) {
		t.Fatalf("OpFail message constant index out of bounds: %d", msgIdx)
	}

	msg, ok := prog.Constants[msgIdx].(runtime.String)
	if !ok {
		t.Fatalf("expected OpFail message constant to be string, got %T", prog.Constants[msgIdx])
	}

	if msg != expectedMessage {
		t.Fatalf("unexpected OpFail message: got %q, want %q", msg, expectedMessage)
	}
}

func TestQueryModifierLowering_ValueUsesLoadNoneAndFail(t *testing.T) {
	prog := compileWithLevel(t, compiler.O0, `RETURN QUERY VALUE ".items" IN @doc USING css`)

	assertThreeSlotQueryDescriptor(t, prog.Bytecode)
	assertFailPrelude(t, prog, runtime.NewString("QUERY VALUE expected at least one match"))

	if !hasOpcode(prog.Bytecode, bytecode.OpLoadIndexConst) {
		t.Fatalf("expected OpLoadIndexConst success path for QUERY VALUE")
	}
}

func TestQueryModifierLowering_OneUsesLoadNoneAndFail(t *testing.T) {
	prog := compileWithLevel(t, compiler.O0, `RETURN QUERY ONE ".items" IN @doc USING css`)

	assertThreeSlotQueryDescriptor(t, prog.Bytecode)
	assertFailPrelude(t, prog, runtime.NewString("QUERY ONE expected exactly one match"))

	if !hasOpcode(prog.Bytecode, bytecode.OpLength) {
		t.Fatalf("expected OpLength for QUERY ONE cardinality check")
	}
	if !hasOpcode(prog.Bytecode, bytecode.OpJumpIfEqConst) {
		t.Fatalf("expected OpJumpIfEqConst for QUERY ONE cardinality check")
	}
}

func TestQueryModifierLowering_ExistsCountAny(t *testing.T) {
	cases := []struct {
		name   string
		expr   string
		opcode bytecode.Opcode
	}{
		{name: "exists", expr: `RETURN QUERY EXISTS ".items" IN @doc USING css`, opcode: bytecode.OpExists},
		{name: "count", expr: `RETURN QUERY COUNT ".items" IN @doc USING css`, opcode: bytecode.OpLength},
		{name: "any", expr: `RETURN QUERY ANY ".items" IN @doc USING css`, opcode: bytecode.OpLoadIndexOptionalConst},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			prog := compileWithLevel(t, compiler.O0, tc.expr)

			assertThreeSlotQueryDescriptor(t, prog.Bytecode)

			if !hasOpcode(prog.Bytecode, tc.opcode) {
				t.Fatalf("expected opcode %s for QUERY %s lowering", tc.opcode, tc.name)
			}

			if hasOpcode(prog.Bytecode, bytecode.OpFail) {
				t.Fatalf("did not expect OpFail in QUERY %s lowering", tc.name)
			}
		})
	}
}
