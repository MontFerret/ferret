package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
)

func collectProjectionValueDefOpcode(t *testing.T, prog *bytecode.Program) bytecode.Opcode {
	t.Helper()

	pushKVIndex, ok := findFirstOpcodeIndex(prog.Bytecode, bytecode.OpPushKV)
	if !ok {
		t.Fatalf("expected OpPushKV in bytecode")
	}

	valueReg := prog.Bytecode[pushKVIndex].Operands[2].Register()
	valueDef, ok := lastRegisterDefOpcodeBefore(prog.Bytecode, pushKVIndex, valueReg)
	if !ok {
		t.Fatalf("expected defining opcode for collect projection register R%d", valueReg)
	}

	return valueDef
}

func TestCollectProjectionO0UsesPlainMoveForAllVarsScopePacking(t *testing.T) {
	prog := compileWithLevel(t, compiler.O0, `
LET users = [{ age: 1 }]

FOR i IN users
	COLLECT g = i.age INTO groups
	RETURN groups
`)

	if got := collectProjectionValueDefOpcode(t, prog); got != bytecode.OpMove {
		t.Fatalf("expected all-vars collect projection handoff to use MOVE, got %s", got.String())
	}
}

func TestCollectProjectionO0UsesPlainMoveForKeepProjectionObject(t *testing.T) {
	prog := compileWithLevel(t, compiler.O0, `
LET users = [{ age: 1 }]

FOR i IN users
	COLLECT g = i.age INTO groups KEEP i
	RETURN groups
`)

	if got := collectProjectionValueDefOpcode(t, prog); got != bytecode.OpMove {
		t.Fatalf("expected KEEP collect projection handoff to use MOVE, got %s", got.String())
	}
}

func TestCollectProjectionO0UsesPlainMoveForTypedCustomProjection(t *testing.T) {
	prog := compileWithLevel(t, compiler.O0, `
LET users = ["alice"]

FOR i IN users
	COLLECT g = i INTO groups = i + "1"
	RETURN groups
`)

	if got := collectProjectionValueDefOpcode(t, prog); got != bytecode.OpMove {
		t.Fatalf("expected typed custom collect projection handoff to use MOVE, got %s", got.String())
	}
}

func TestCollectProjectionO0KeepsTrackedMoveForUnknownCustomProjection(t *testing.T) {
	prog := compileWithLevel(t, compiler.O0, `
FUNC project(v) => v

LET users = [1]

FOR i IN users
	COLLECT g = i INTO groups = project(i)
	RETURN groups
`)

	if got := collectProjectionValueDefOpcode(t, prog); got != bytecode.OpMoveTracked {
		t.Fatalf("expected unknown custom collect projection handoff to use MOVET, got %s", got.String())
	}
}
