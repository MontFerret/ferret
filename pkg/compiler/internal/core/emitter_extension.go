package core

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
)

// Emitter Helpers - Common opcode shortcut methods

// ─── Loop & State ──────────────────────────────────────────────────────

func (e *Emitter) EmitIter(dst, src bytecode.Operand) {
	e.EmitAB(bytecode.OpIter, dst, src)
}

func (e *Emitter) EmitIterNext(iterator bytecode.Operand, label Label) {
	e.EmitJumpc(bytecode.OpIterNext, iterator, label)
}

func (e *Emitter) EmitIterKey(dst, iterator bytecode.Operand) {
	e.EmitAB(bytecode.OpIterKey, dst, iterator)
}

func (e *Emitter) EmitIterValue(dst, iterator bytecode.Operand) {
	e.EmitAB(bytecode.OpIterValue, dst, iterator)
}

func (e *Emitter) EmitIterSkip(state, count bytecode.Operand, label Label) {
	e.emitInstruction(bytecode.Instruction{
		Opcode:   bytecode.OpIterSkip,
		Operands: [3]bytecode.Operand{jumpPlaceholder, state, count},
	})

	pos := len(e.instructions) - 1
	e.addLabelRef(pos, 0, label)
}

func (e *Emitter) EmitIterLimit(state, count bytecode.Operand, label Label) {
	e.emitInstruction(bytecode.Instruction{
		Opcode:   bytecode.OpIterLimit,
		Operands: [3]bytecode.Operand{jumpPlaceholder, state, count},
	})

	pos := len(e.instructions) - 1
	e.addLabelRef(pos, 0, label)
}

// ─── Value & Memory ──────────────────────────────────────────────────────

func (e *Emitter) EmitMove(dst, src bytecode.Operand) {
	e.EmitAB(bytecode.OpMoveTracked, dst, src)
}

func (e *Emitter) EmitPlainMove(dst, src bytecode.Operand) {
	e.EmitAB(bytecode.OpMove, dst, src)
}

func (e *Emitter) EmitMoveTracked(dst, src bytecode.Operand) {
	e.EmitAB(bytecode.OpMoveTracked, dst, src)
}

func (e *Emitter) EmitPush(dst, src bytecode.Operand) {
	e.EmitAB(bytecode.OpPush, dst, src)
}

func (e *Emitter) EmitArrayPush(dst, src bytecode.Operand) {
	e.EmitAB(bytecode.OpArrayPush, dst, src)
}

func (e *Emitter) EmitPushKV(dst, key, val bytecode.Operand) {
	e.EmitABC(bytecode.OpPushKV, dst, key, val)
}

func (e *Emitter) EmitCounterInc(dst bytecode.Operand) {
	e.EmitA(bytecode.OpCounterInc, dst)
}

func (e *Emitter) EmitObjectSet(dst, key, val bytecode.Operand) {
	e.EmitABC(bytecode.OpObjectSet, dst, key, val)
}

func (e *Emitter) EmitObjectSetConst(dst, keyConst, val bytecode.Operand) {
	e.EmitABC(bytecode.OpObjectSetConst, dst, keyConst, val)
}

func (e *Emitter) EmitClose(reg bytecode.Operand) {
	e.EmitA(bytecode.OpClose, reg)
}

func (e *Emitter) EmitLoadNone(dst bytecode.Operand) {
	e.EmitA(bytecode.OpLoadNone, dst)
}

func (e *Emitter) EmitLoadConst(dst bytecode.Operand, constant bytecode.Operand) {
	e.EmitAB(bytecode.OpLoadConst, dst, constant)
}

func (e *Emitter) EmitLoadParam(dst, slot bytecode.Operand) {
	e.EmitAB(bytecode.OpLoadParam, dst, slot)
}

func (e *Emitter) EmitLoadAggregateKey(dst, key, selector bytecode.Operand) {
	e.EmitABC(bytecode.OpLoadAggregateKey, dst, key, selector)
}

func (e *Emitter) EmitAggregateUpdate(collector, value bytecode.Operand, selector int) {
	e.emitInstructionWithSelectorSlot(bytecode.Instruction{
		Opcode:   bytecode.OpAggregateUpdate,
		Operands: [3]bytecode.Operand{collector, value, bytecode.NoopOperand},
	}, selector)
}

func (e *Emitter) EmitAggregateGroupUpdate(collector, key, value bytecode.Operand, selector int) {
	e.emitInstructionWithSelectorSlot(bytecode.Instruction{
		Opcode:   bytecode.OpAggregateGroupUpdate,
		Operands: [3]bytecode.Operand{collector, key, value},
	}, selector)
}

func (e *Emitter) EmitBoolean(dst bytecode.Operand, value bool) {
	if value {
		e.EmitAB(bytecode.OpLoadBool, dst, 1)
	} else {
		e.EmitAB(bytecode.OpLoadBool, dst, 0)
	}
}

// ─── Data Structures ──────────────────────────────────────────────────────

func (e *Emitter) EmitArray(dst bytecode.Operand, size int) {
	e.EmitAB(bytecode.OpLoadArray, dst, bytecode.Operand(size))
}

func (e *Emitter) EmitObject(dst bytecode.Operand, size int) {
	e.EmitAB(bytecode.OpLoadObject, dst, bytecode.Operand(size))
}

func (e *Emitter) EmitRange(dst, start, end bytecode.Operand) {
	e.EmitABC(bytecode.OpLoadRange, dst, start, end)
}

func (e *Emitter) EmitLoadIndex(dst, arr, idx bytecode.Operand) {
	e.EmitABC(bytecode.OpLoadIndex, dst, arr, idx)
}

func (e *Emitter) EmitLoadKey(dst, obj, key bytecode.Operand) {
	e.EmitABC(bytecode.OpLoadKey, dst, obj, key)
}

func (e *Emitter) EmitLoadProperty(dst, obj, prop bytecode.Operand) {
	e.EmitABC(bytecode.OpLoadProperty, dst, obj, prop)
}

func (e *Emitter) EmitLoadPropertyOptional(dst, obj, prop bytecode.Operand) {
	e.EmitABC(bytecode.OpLoadPropertyOptional, dst, obj, prop)
}

// ─── Arithmetic and Logical ──────────────────────────────────────────────

func (e *Emitter) EmitAdd(dst, a, b bytecode.Operand) {
	e.EmitABC(bytecode.OpAdd, dst, a, b)
}

func (e *Emitter) EmitSub(dst, a, b bytecode.Operand) {
	e.EmitABC(bytecode.OpSub, dst, a, b)
}

func (e *Emitter) EmitMul(dst, a, b bytecode.Operand) {
	e.EmitABC(bytecode.OpMul, dst, a, b)
}

func (e *Emitter) EmitDiv(dst, a, b bytecode.Operand) {
	e.EmitABC(bytecode.OpDiv, dst, a, b)
}

func (e *Emitter) EmitMod(dst, a, b bytecode.Operand) {
	e.EmitABC(bytecode.OpMod, dst, a, b)
}

func (e *Emitter) EmitEq(dst, a, b bytecode.Operand) {
	e.EmitABC(bytecode.OpEq, dst, a, b)
}

func (e *Emitter) EmitNeq(dst, a, b bytecode.Operand) {
	e.EmitABC(bytecode.OpNe, dst, a, b)
}

func (e *Emitter) EmitGt(dst, a, b bytecode.Operand) {
	e.EmitABC(bytecode.OpGt, dst, a, b)
}

func (e *Emitter) EmitLt(dst, a, b bytecode.Operand) {
	e.EmitABC(bytecode.OpLt, dst, a, b)
}

func (e *Emitter) EmitGte(dst, a, b bytecode.Operand) {
	e.EmitABC(bytecode.OpGte, dst, a, b)
}

func (e *Emitter) EmitLte(dst, a, b bytecode.Operand) {
	e.EmitABC(bytecode.OpLte, dst, a, b)
}

// ─── Control Flow ────────────────────────────────────────────────────────

func (e *Emitter) EmitJump(label Label) {
	e.EmitA(bytecode.OpJump, bytecode.Operand(jumpPlaceholder))
	pos := len(e.instructions) - 1
	e.addLabelRef(pos, 0, label)
}

// EmitJumpAB emits a jump opcode with a state and an argument.
func (e *Emitter) EmitJumpAB(op bytecode.Opcode, state, cond bytecode.Operand, label Label) {
	e.EmitABC(op, state, cond, bytecode.Operand(jumpPlaceholder))
	pos := len(e.instructions) - 1
	e.addLabelRef(pos, 2, label)
}

// EmitJumpc emits a conditional jump opcode.
func (e *Emitter) EmitJumpc(op bytecode.Opcode, reg bytecode.Operand, label Label) {
	e.EmitAB(op, bytecode.Operand(jumpPlaceholder), reg)
	pos := len(e.instructions) - 1
	e.addLabelRef(pos, 0, label)
}

// EmitJumpCompare emits a conditional jump that compares two operands.
func (e *Emitter) EmitJumpCompare(op bytecode.Opcode, left, right bytecode.Operand, label Label) {
	e.emitInstruction(bytecode.Instruction{
		Opcode:   op,
		Operands: [3]bytecode.Operand{jumpPlaceholder, left, right},
	})
	pos := len(e.instructions) - 1
	e.addLabelRef(pos, 0, label)
}

func (e *Emitter) EmitJumpIfFalse(cond bytecode.Operand, label Label) {
	e.EmitJumpIf(cond, false, label)
}

func (e *Emitter) EmitJumpIfTrue(cond bytecode.Operand, label Label) {
	e.EmitJumpIf(cond, true, label)
}

func (e *Emitter) EmitJumpIfNone(cond bytecode.Operand, label Label) {
	e.EmitJumpc(bytecode.OpJumpIfNone, cond, label)
}

func (e *Emitter) EmitJumpIf(cond bytecode.Operand, isTrue bool, label Label) {
	if isTrue {
		e.EmitJumpc(bytecode.OpJumpIfTrue, cond, label)
		return
	}

	e.EmitJumpc(bytecode.OpJumpIfFalse, cond, label)
}
