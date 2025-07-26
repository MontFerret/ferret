package core

import (
	"github.com/MontFerret/ferret/pkg/vm"
)

// Emitter Helpers - Common opcode shortcut methods

// ─── Loop & Iterator ──────────────────────────────────────────────────────

func (e *Emitter) EmitIter(dst, src vm.Operand) {
	e.EmitAB(vm.OpIter, dst, src)
}

func (e *Emitter) EmitIterNext(iterator vm.Operand, label Label) {
	e.EmitJumpc(vm.OpIterNext, iterator, label)
}

func (e *Emitter) EmitIterKey(dst, iterator vm.Operand) {
	e.EmitAB(vm.OpIterKey, dst, iterator)
}

func (e *Emitter) EmitIterValue(dst, iterator vm.Operand) {
	e.EmitAB(vm.OpIterValue, dst, iterator)
}

func (e *Emitter) EmitIterSkip(state, count vm.Operand, label Label) {
	e.instructions = append(e.instructions, vm.Instruction{
		Opcode:   vm.OpIterSkip,
		Operands: [3]vm.Operand{jumpPlaceholder, state, count},
	})

	pos := len(e.instructions) - 1
	e.addLabelRef(pos, 0, label)
}

func (e *Emitter) EmitIterLimit(state, count vm.Operand, label Label) {
	e.instructions = append(e.instructions, vm.Instruction{
		Opcode:   vm.OpIterLimit,
		Operands: [3]vm.Operand{jumpPlaceholder, state, count},
	})

	pos := len(e.instructions) - 1
	e.addLabelRef(pos, 0, label)
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

func (e *Emitter) EmitLoadNone(dst vm.Operand) {
	e.EmitA(vm.OpLoadNone, dst)
}

func (e *Emitter) EmitLoadConst(dst vm.Operand, constant vm.Operand) {
	e.EmitAB(vm.OpLoadConst, dst, constant)
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

func (e *Emitter) EmitArray(dst vm.Operand, seq RegisterSequence) {
	if len(seq) > 0 {
		e.EmitAs(vm.OpLoadArray, dst, seq)
	} else {
		e.EmitA(vm.OpLoadArray, dst)
	}
}

func (e *Emitter) EmitObject(dst vm.Operand, seq RegisterSequence) {
	if len(seq) > 0 {
		e.EmitAs(vm.OpLoadObject, dst, seq)
	} else {
		e.EmitA(vm.OpLoadObject, dst)
	}
}

func (e *Emitter) EmitRange(dst, start, end vm.Operand) {
	e.EmitABC(vm.OpLoadRange, dst, start, end)
}

func (e *Emitter) EmitLoadIndex(dst, arr, idx vm.Operand) {
	e.EmitABC(vm.OpLoadIndex, dst, arr, idx)
}

func (e *Emitter) EmitLoadKey(dst, obj, key vm.Operand) {
	e.EmitABC(vm.OpLoadKey, dst, obj, key)
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

func (e *Emitter) EmitJump(label Label) {
	e.EmitA(vm.OpJump, vm.Operand(jumpPlaceholder))
	pos := len(e.instructions) - 1
	e.addLabelRef(pos, 0, label)
}

// EmitJumpAB emits a jump opcode with a state and an argument.
func (e *Emitter) EmitJumpAB(op vm.Opcode, state, cond vm.Operand, label Label) {
	e.EmitABC(op, state, cond, vm.Operand(jumpPlaceholder))
	pos := len(e.instructions) - 1
	e.addLabelRef(pos, 2, label)
}

// EmitJumpc emits a conditional jump opcode.
func (e *Emitter) EmitJumpc(op vm.Opcode, reg vm.Operand, label Label) {
	e.EmitAB(op, vm.Operand(jumpPlaceholder), reg)
	pos := len(e.instructions) - 1
	e.addLabelRef(pos, 0, label)
}

func (e *Emitter) EmitJumpIfFalse(cond vm.Operand, label Label) {
	e.EmitJumpIf(cond, false, label)
}

func (e *Emitter) EmitJumpIfTrue(cond vm.Operand, label Label) {
	e.EmitJumpIf(cond, true, label)
}

func (e *Emitter) EmitJumpIf(cond vm.Operand, isTrue bool, label Label) {
	if isTrue {
		e.EmitJumpc(vm.OpJumpIfTrue, cond, label)
		return
	}

	e.EmitJumpc(vm.OpJumpIfFalse, cond, label)
}
