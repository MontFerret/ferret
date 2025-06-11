package core

import (
	"github.com/MontFerret/ferret/pkg/vm"
)

// Emitter Helpers - Common opcode shortcut methods

// ─── Loop & Iterator ──────────────────────────────────────────────────────

func (e *Emitter) EmitIter(dst, src vm.Operand) {
	e.EmitAB(vm.OpIter, dst, src)
}

func (e *Emitter) EmitIterNext(jumpTarget int, iterator vm.Operand) int {
	return e.EmitJumpc(vm.OpIterNext, jumpTarget, iterator)
}

func (e *Emitter) EmitIterKey(dst, iterator vm.Operand) {
	e.EmitAB(vm.OpIterKey, dst, iterator)
}

func (e *Emitter) EmitIterValue(dst, iterator vm.Operand) {
	e.EmitAB(vm.OpIterValue, dst, iterator)
}

func (e *Emitter) EmitIterSkip(state, count vm.Operand, jump int) {
	e.EmitABx(vm.OpIterSkip, state, count, jump)
}

func (e *Emitter) EmitIterLimit(state, count vm.Operand, jump int) {
	e.EmitABx(vm.OpIterLimit, state, count, jump)
}

// ─── Value & Memory ──────────────────────────────────────────────────────

func (e *Emitter) EmitMove(dst, src vm.Operand) {
	e.EmitAB(vm.OpMove, dst, src)
}

func (e *Emitter) EmitPush(dst, src vm.Operand) {
	e.EmitAB(vm.OpPush, dst, src)
}

func (e *Emitter) EmitPushKV(dst, key, val vm.Operand) {
	e.EmitABC(vm.OpPushKV, dst, key, val)
}

func (e *Emitter) EmitClose(reg vm.Operand) {
	e.EmitA(vm.OpClose, reg)
}

func (e *Emitter) EmitLoadConst(dst vm.Operand, constant vm.Operand) {
	e.EmitAB(vm.OpLoadConst, dst, constant)
}

func (e *Emitter) EmitLoadGlobal(dst, constant vm.Operand) {
	e.EmitAB(vm.OpLoadGlobal, dst, constant)
}

func (e *Emitter) EmitLoadParam(dst, constant vm.Operand) {
	e.EmitAB(vm.OpLoadParam, dst, constant)
}

func (e *Emitter) EmitBoolean(dst vm.Operand, value bool) {
	if value {
		e.EmitAB(vm.OpLoadBool, dst, 1)
	} else {
		e.EmitAB(vm.OpLoadBool, dst, 0)
	}
}

// ─── Data Structures ──────────────────────────────────────────────────────

func (e *Emitter) EmitEmptyList(dst vm.Operand) {
	e.EmitA(vm.OpList, dst)
}

func (e *Emitter) EmitList(dst vm.Operand, seq RegisterSequence) {
	e.EmitAs(vm.OpList, dst, seq)
}

func (e *Emitter) EmitEmptyMap(dst vm.Operand) {
	e.EmitA(vm.OpMap, dst)
}

func (e *Emitter) EmitMap(dst vm.Operand, seq RegisterSequence) {
	e.EmitAs(vm.OpMap, dst, seq)
}

func (e *Emitter) EmitRange(dst, start, end vm.Operand) {
	e.EmitABC(vm.OpRange, dst, start, end)
}

func (e *Emitter) EmitLoadIndex(dst, arr, idx vm.Operand) {
	e.EmitABC(vm.OpLoadIndex, dst, arr, idx)
}

func (e *Emitter) EmitLoadProperty(dst, obj, prop vm.Operand) {
	e.EmitABC(vm.OpLoadProperty, dst, obj, prop)
}

func (e *Emitter) EmitLoadPropertyOptional(dst, obj, prop vm.Operand) {
	e.EmitABC(vm.OpLoadPropertyOptional, dst, obj, prop)
}

// ─── Arithmetic and Logical ──────────────────────────────────────────────

func (e *Emitter) EmitAdd(dst, a, b vm.Operand) {
	e.EmitABC(vm.OpAdd, dst, a, b)
}

func (e *Emitter) EmitSub(dst, a, b vm.Operand) {
	e.EmitABC(vm.OpSub, dst, a, b)
}

func (e *Emitter) EmitMul(dst, a, b vm.Operand) {
	e.EmitABC(vm.OpMulti, dst, a, b)
}

func (e *Emitter) EmitDiv(dst, a, b vm.Operand) {
	e.EmitABC(vm.OpDiv, dst, a, b)
}

func (e *Emitter) EmitMod(dst, a, b vm.Operand) {
	e.EmitABC(vm.OpMod, dst, a, b)
}

func (e *Emitter) EmitEq(dst, a, b vm.Operand) {
	e.EmitABC(vm.OpEq, dst, a, b)
}

func (e *Emitter) EmitNeq(dst, a, b vm.Operand) {
	e.EmitABC(vm.OpNeq, dst, a, b)
}

func (e *Emitter) EmitGt(dst, a, b vm.Operand) {
	e.EmitABC(vm.OpGt, dst, a, b)
}

func (e *Emitter) EmitLt(dst, a, b vm.Operand) {
	e.EmitABC(vm.OpLt, dst, a, b)
}

func (e *Emitter) EmitGte(dst, a, b vm.Operand) {
	e.EmitABC(vm.OpGte, dst, a, b)
}

func (e *Emitter) EmitLte(dst, a, b vm.Operand) {
	e.EmitABC(vm.OpLte, dst, a, b)
}

// ─── Control Flow ────────────────────────────────────────────────────────

func (e *Emitter) EmitJump(pos int) int {
	e.EmitA(vm.OpJump, vm.Operand(pos))

	return len(e.instructions) - 1
}

// EmitJumpAB emits a jump opcode with a state and an argument.
func (e *Emitter) EmitJumpAB(op vm.Opcode, state, cond vm.Operand, pos int) int {
	e.EmitABC(op, state, cond, vm.Operand(pos))

	return len(e.instructions) - 1
}

// EmitJumpc emits a conditional jump opcode.
func (e *Emitter) EmitJumpc(op vm.Opcode, pos int, reg vm.Operand) int {
	e.EmitAB(op, vm.Operand(pos), reg)

	return len(e.instructions) - 1
}

func (e *Emitter) EmitJumpIfFalse(cond vm.Operand, jumpTarget int) int {
	return e.EmitJumpIf(cond, false, jumpTarget)
}

func (e *Emitter) EmitJumpIfTrue(cond vm.Operand, jumpTarget int) int {
	return e.EmitJumpIf(cond, true, jumpTarget)
}

func (e *Emitter) EmitJumpIf(cond vm.Operand, isTrue bool, jumpTarget int) int {
	if isTrue {
		return e.EmitJumpc(vm.OpJumpIfTrue, jumpTarget, cond)
	}

	return e.EmitJumpc(vm.OpJumpIfFalse, jumpTarget, cond)
}

func (e *Emitter) EmitReturnValue(val vm.Operand) {
	if val.IsConstant() {
		e.EmitAB(vm.OpLoadGlobal, vm.NoopOperand, val)
	} else {
		e.EmitAB(vm.OpMove, vm.NoopOperand, val)
	}

	e.Emit(vm.OpReturn)
}
