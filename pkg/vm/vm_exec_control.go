package vm

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/operators"
)

func (vm *VM) execControlOps(
	ctx context.Context,
	op bytecode.Opcode,
	dst, src1, src2 bytecode.Operand,
	reg []runtime.Value,
	constants []runtime.Value,
) (bool, error) {
	switch op {
	case bytecode.OpJumpIfNe:
		if operators.NotEquals(ctx, reg[src1], reg[src2]) {
			vm.pc = int(dst)
		}
	case bytecode.OpJumpIfNeConst:
		if operators.NotEquals(ctx, reg[src1], constants[src2.Constant()]) {
			vm.pc = int(dst)
		}
	case bytecode.OpJumpIfEq:
		if operators.Equals(ctx, reg[src1], reg[src2]) {
			vm.pc = int(dst)
		}
	case bytecode.OpJumpIfEqConst:
		if operators.Equals(ctx, reg[src1], constants[src2.Constant()]) {
			vm.pc = int(dst)
		}
	case bytecode.OpJumpIfMissingProperty:
		obj, ok := reg[src1].(runtime.Map)
		if !ok {
			vm.pc = int(dst)
			return false, nil
		}

		key, ok := reg[src2].(runtime.String)
		if !ok {
			vm.pc = int(dst)
			return false, nil
		}

		has, err := obj.ContainsKey(ctx, key)
		if err != nil {
			return false, vm.handleProtectedError(err)
		}
		if !has {
			vm.pc = int(dst)
		}
	case bytecode.OpJumpIfMissingPropertyConst:
		obj, ok := reg[src1].(runtime.Map)
		if !ok {
			vm.pc = int(dst)
			return false, nil
		}

		key, ok := constants[src2.Constant()].(runtime.String)
		if !ok {
			vm.pc = int(dst)
			return false, nil
		}

		has, err := obj.ContainsKey(ctx, key)
		if err != nil {
			return false, vm.handleProtectedError(err)
		}
		if !has {
			vm.pc = int(dst)
		}
	case bytecode.OpReturn:
		retVal := reg[dst]

		if frame, ok := vm.popFrame(); ok {
			vm.regPool.put(vm.registers.Values)
			vm.registers.Values = frame.registers
			vm.registers.Values[frame.returnDest] = retVal
			vm.pc = frame.returnPC
			return false, nil
		}

		reg[bytecode.NoopOperand] = retVal
		return true, nil
	}

	return false, nil
}
