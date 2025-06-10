package internal

import (
	"github.com/MontFerret/ferret/pkg/runtime"
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

func (e *Emitter) EmitLoadConst(dst vm.Operand, val runtime.Value, symbols *SymbolTable) {
	e.EmitAB(vm.OpLoadConst, dst, symbols.AddConstant(val))
}

// ─── Data Structures ──────────────────────────────────────────────────────

func (e *Emitter) EmitList(dst vm.Operand, seq RegisterSequence) {
	e.EmitAs(vm.OpList, dst, seq)
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
