package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
)

func TestSort(t *testing.T) {
	RunUseCases(t, []UseCase{
		SkipByteCodeCase(`
FOR s IN []
	SORT s
	RETURN s
`, BC{
			I(bytecode.OpReturn, 0, 7),
		}),
	})
}

func findNthOpcodeIndex(code []bytecode.Instruction, op bytecode.Opcode, nth int) (int, bool) {
	count := 0

	for i, inst := range code {
		if inst.Opcode != op {
			continue
		}

		if count == nth {
			return i, true
		}

		count++
	}

	return -1, false
}

func TestForSortO0UsesPlainMovesForKeyAndScopeProjection(t *testing.T) {
	prog := compileWithLevel(t, compiler.O0, `
LET strs = ["foo", "bar", "qaz", "abc"]

FOR s IN strs
	SORT s + "1"
	RETURN s
`)

	pushKVIndex, ok := findFirstOpcodeIndex(prog.Bytecode, bytecode.OpPushKV)
	if !ok {
		t.Fatalf("expected OpPushKV in bytecode")
	}

	keyReg := prog.Bytecode[pushKVIndex].Operands[1].Register()
	keyDef, ok := lastRegisterDefOpcodeBefore(prog.Bytecode, pushKVIndex, keyReg)
	if !ok {
		t.Fatalf("expected defining opcode for sort key register R%d", keyReg)
	}
	if keyDef != bytecode.OpMove {
		t.Fatalf("expected sort key copy to use MOVE, got %s", keyDef.String())
	}

	valueReg := prog.Bytecode[pushKVIndex].Operands[2].Register()
	valueDef, ok := lastRegisterDefOpcodeBefore(prog.Bytecode, pushKVIndex, valueReg)
	if !ok {
		t.Fatalf("expected defining opcode for projected scope register R%d", valueReg)
	}
	if valueDef != bytecode.OpMove {
		t.Fatalf("expected projected scope handoff to use MOVE, got %s", valueDef.String())
	}
}

func TestForSortO1UsesPlainMoveForScopeProjectionAndTrackedSorterTransfer(t *testing.T) {
	prog := compileWithLevel(t, compiler.O1, `
LET strs = ["foo", "bar", "qaz", "abc"]

FOR s IN strs
	SORT s + "1"
	RETURN s
`)

	pushKVIndex, ok := findFirstOpcodeIndex(prog.Bytecode, bytecode.OpPushKV)
	if !ok {
		t.Fatalf("expected OpPushKV in bytecode")
	}

	valueReg := prog.Bytecode[pushKVIndex].Operands[2].Register()
	valueDef, ok := lastRegisterDefOpcodeBefore(prog.Bytecode, pushKVIndex, valueReg)
	if !ok {
		t.Fatalf("expected defining opcode for projected scope register R%d", valueReg)
	}
	if valueDef != bytecode.OpMove {
		t.Fatalf("expected projected scope handoff to use MOVE, got %s", valueDef.String())
	}

	secondIterIndex, ok := findNthOpcodeIndex(prog.Bytecode, bytecode.OpIter, 1)
	if !ok {
		t.Fatalf("expected second OpIter in bytecode")
	}

	sortedSrcReg := prog.Bytecode[secondIterIndex].Operands[1].Register()
	sortedSrcDef, ok := lastRegisterDefOpcodeBefore(prog.Bytecode, secondIterIndex, sortedSrcReg)
	if !ok {
		t.Fatalf("expected defining opcode for sorted source register R%d", sortedSrcReg)
	}
	if sortedSrcDef != bytecode.OpMoveTracked {
		t.Fatalf("expected sorter transfer to remain MOVET, got %s", sortedSrcDef.String())
	}
}
