package optimization

import (
	"context"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func foldInstructionByOpcode(inst *bytecode.Instruction, env constFoldEnv) constFoldResult {
	result := newConstFoldResult()

	switch {
	case foldLoadInstruction(inst, env, &result):
	case foldMoveInstruction(inst, env, &result):
	case foldIncDecInstruction(inst, env, &result):
	case foldUnaryInstruction(inst, env, &result):
	case foldBinaryInstruction(inst, env, &result):
	case foldConcatInstruction(inst, env, &result):
	}

	return result
}

func foldLoadInstruction(inst *bytecode.Instruction, env constFoldEnv, result *constFoldResult) bool {
	dst, ok := registerOperand(inst.Operands[0])
	if !ok {
		return false
	}

	switch inst.Opcode {
	case bytecode.OpLoadConst:
		val := env.program.Constants[inst.Operands[1].Constant()]
		if isSimpleConst(val) {
			result.setConst(dst, val)
		}
		return true
	case bytecode.OpLoadNone:
		result.setConst(dst, runtime.None)
		return true
	case bytecode.OpLoadZero:
		result.setConst(dst, runtime.ZeroInt)
		return true
	case bytecode.OpLoadBool:
		result.setConst(dst, runtime.Boolean(inst.Operands[1] == 1))
		return true
	default:
		return false
	}
}

func foldMoveInstruction(inst *bytecode.Instruction, env constFoldEnv, result *constFoldResult) bool {
	if inst.Opcode != bytecode.OpMove && inst.Opcode != bytecode.OpMoveTracked {
		return false
	}

	dst, ok := registerOperand(inst.Operands[0])
	if !ok || !inst.Operands[1].IsRegister() {
		return true
	}

	val, ok := env.state[inst.Operands[1].Register()]
	if !ok {
		return true
	}

	result.rewriteWithConst(inst, dst, val, env)

	return true
}

func foldIncDecInstruction(inst *bytecode.Instruction, env constFoldEnv, result *constFoldResult) bool {
	if inst.Opcode != bytecode.OpIncr && inst.Opcode != bytecode.OpDecr {
		return false
	}

	dst, ok := registerOperand(inst.Operands[0])
	if !ok {
		return true
	}

	val, ok := env.state[dst]
	if !ok {
		return true
	}

	out, ok := foldIncDecValue(inst.Opcode, val, env.bg)
	if !ok || !isSimpleConst(out) {
		return true
	}

	result.rewriteWithConst(inst, dst, out, env)

	return true
}

func foldIncDecValue(op bytecode.Opcode, val runtime.Value, bg context.Context) (runtime.Value, bool) {
	if op == bytecode.OpIncr {
		return increment(bg, val), true
	}

	if op == bytecode.OpDecr {
		return decrement(bg, val), true
	}

	return nil, false
}

func foldUnaryInstruction(inst *bytecode.Instruction, env constFoldEnv, result *constFoldResult) bool {
	if !isUnaryFoldOpcode(inst.Opcode) {
		return false
	}

	dst, ok := registerOperand(inst.Operands[0])
	if !ok || !inst.Operands[1].IsRegister() {
		return true
	}

	val, ok := env.state[inst.Operands[1].Register()]
	if !ok {
		return true
	}

	out, ok := foldUnary(inst.Opcode, val, env.bg)
	if !ok || !isSimpleConst(out) {
		return true
	}

	result.rewriteWithConst(inst, dst, out, env)

	return true
}

func isUnaryFoldOpcode(op bytecode.Opcode) bool {
	switch op {
	case bytecode.OpCastBool, bytecode.OpNegate, bytecode.OpFlipPositive, bytecode.OpFlipNegative, bytecode.OpNot:
		return true
	default:
		return false
	}
}

func foldBinaryInstruction(inst *bytecode.Instruction, env constFoldEnv, result *constFoldResult) bool {
	if !isBinaryFoldOpcode(inst.Opcode) {
		return false
	}

	dst, ok := registerOperand(inst.Operands[0])
	if !ok {
		return true
	}

	left, right, ok := resolveBinaryFoldOperands(*inst, env.state, env.program)
	if !ok {
		return true
	}

	foldOpcode := normalizeBinaryFoldOpcode(inst.Opcode)
	out, ok := foldBinary(foldOpcode, left, right, env.bg)
	if !ok || !isSimpleConst(out) {
		return true
	}

	result.rewriteWithConst(inst, dst, out, env)

	return true
}

func isBinaryFoldOpcode(op bytecode.Opcode) bool {
	switch op {
	case bytecode.OpAdd, bytecode.OpAddConst, bytecode.OpSub, bytecode.OpMul, bytecode.OpDiv, bytecode.OpMod,
		bytecode.OpCmp, bytecode.OpEq, bytecode.OpNe, bytecode.OpGt, bytecode.OpLt, bytecode.OpGte, bytecode.OpLte:
		return true
	default:
		return false
	}
}

func normalizeBinaryFoldOpcode(op bytecode.Opcode) bytecode.Opcode {
	if op == bytecode.OpAddConst {
		return bytecode.OpAdd
	}

	return op
}

func resolveBinaryFoldOperands(inst bytecode.Instruction, state constState, program *bytecode.Program) (runtime.Value, runtime.Value, bool) {
	if !inst.Operands[1].IsRegister() {
		return nil, nil, false
	}

	left, ok := state[inst.Operands[1].Register()]
	if !ok {
		return nil, nil, false
	}

	if inst.Opcode == bytecode.OpAddConst {
		if !inst.Operands[2].IsConstant() {
			return nil, nil, false
		}

		return left, program.Constants[inst.Operands[2].Constant()], true
	}

	if !inst.Operands[2].IsRegister() {
		return nil, nil, false
	}

	right, ok := state[inst.Operands[2].Register()]
	if !ok {
		return nil, nil, false
	}

	return left, right, true
}

func foldConcatInstruction(inst *bytecode.Instruction, env constFoldEnv, result *constFoldResult) bool {
	if inst.Opcode != bytecode.OpConcat {
		return false
	}

	dst, start, count, ok := resolveConcatFoldOperands(*inst)
	if !ok {
		return true
	}

	if count <= 0 {
		result.rewriteWithConst(inst, dst, runtime.EmptyString, env)
		return true
	}

	if start <= 0 {
		return true
	}

	out, ok := buildConcatFoldConst(env.state, start, count)
	if !ok {
		return true
	}

	result.rewriteWithConst(inst, dst, out, env)

	return true
}

func resolveConcatFoldOperands(inst bytecode.Instruction) (dst int, start int, count int, ok bool) {
	if !inst.Operands[0].IsRegister() || !inst.Operands[1].IsRegister() {
		return 0, 0, 0, false
	}

	return inst.Operands[0].Register(), inst.Operands[1].Register(), int(inst.Operands[2]), true
}

func buildConcatFoldConst(state constState, start, count int) (runtime.Value, bool) {
	var builder strings.Builder

	for i := 0; i < count; i++ {
		val, ok := state[start+i]
		if !ok {
			return nil, false
		}

		builder.WriteString(runtime.ToString(val).String())
	}

	return runtime.NewString(builder.String()), true
}

func registerOperand(op bytecode.Operand) (int, bool) {
	if !op.IsRegister() {
		return 0, false
	}

	return op.Register(), true
}
