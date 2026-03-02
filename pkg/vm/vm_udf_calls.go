package vm

import (
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func (vm *VM) resolveUdfID(val runtime.Value) (int, error) {
	idVal, ok := val.(runtime.Int)
	if !ok {
		return -1, ErrInvalidFunctionName
	}

	return int(idVal), nil
}

func (vm *VM) udfByID(id int) (*bytecode.UDF, error) {
	if id < 0 || vm.program == nil || id >= len(vm.program.Functions.UserDefined) {
		return nil, ErrUnresolvedFunction
	}

	return &vm.program.Functions.UserDefined[id], nil
}

func udfArgRange(src1, src2 bytecode.Operand) (int, int, bool) {
	if !src1.IsRegister() || !src2.IsRegister() {
		return 0, 0, false
	}

	start := src1.Register()
	end := src2.Register()

	if start <= 0 || end < start {
		return 0, 0, false
	}

	return start, end, true
}

func udfArgCount(src1, src2 bytecode.Operand) int {
	start, end, ok := udfArgRange(src1, src2)
	if !ok {
		return 0
	}

	return end - start + 1
}

func copyUdfArgs(dst, src []runtime.Value, src1, src2 bytecode.Operand) {
	if len(dst) <= 1 || len(src) == 0 {
		return
	}

	start, end, ok := udfArgRange(src1, src2)
	if !ok {
		return
	}

	if start >= len(src) {
		return
	}

	count := end - start + 1
	if maxCount := len(dst) - 1; count > maxCount {
		count = maxCount
	}

	if maxCount := len(src) - start; count > maxCount {
		count = maxCount
	}

	copy(dst[1:1+count], src[start:start+count])
}

func collectUdfArgs(src []runtime.Value, src1, src2 bytecode.Operand) []runtime.Value {
	start, end, ok := udfArgRange(src1, src2)
	if !ok || start >= len(src) {
		return nil
	}

	count := end - start + 1
	if maxCount := len(src) - start; count > maxCount {
		count = maxCount
	}

	if count <= 0 {
		return nil
	}

	out := make([]runtime.Value, count)
	copy(out, src[start:start+count])
	return out
}

func (vm *VM) callUdf(op bytecode.Opcode, dst, src1, src2 bytecode.Operand) error {
	reg := vm.registers.Values

	fnID, err := vm.resolveUdfID(reg[dst])
	if err != nil {
		return err
	}

	udf, err := vm.udfByID(fnID)
	if err != nil {
		return err
	}

	argCount := udfArgCount(src1, src2)
	if udf.Params != argCount {
		return runtime.Error(runtime.ErrInvalidArgument, fmt.Sprintf("UDF '%s' expects %d arguments, got %d", udf.Name, udf.Params, argCount))
	}

	if udf.Registers <= 0 {
		return runtime.Error(runtime.ErrInvalidOperation, fmt.Sprintf("UDF '%s' has invalid register window", udf.Name))
	}

	newRegs := vm.regPool.get(udf.Registers)
	copyUdfArgs(newRegs, reg, src1, src2)

	vm.pushFrame(vm.pc, dst, isProtectedUdfCall(op), fnID)
	vm.registers.Values = newRegs
	vm.pc = udf.Entry

	return nil
}

func (vm *VM) tailCallUdf(dst, src1, src2 bytecode.Operand) error {
	if len(vm.frames) == 0 {
		return ErrUnresolvedFunction
	}

	reg := vm.registers.Values
	fnID, err := vm.resolveUdfID(reg[dst])
	if err != nil {
		return err
	}

	udf, err := vm.udfByID(fnID)
	if err != nil {
		return err
	}

	args := collectUdfArgs(reg, src1, src2)
	argCount := len(args)
	if udf.Params != argCount {
		return runtime.Error(runtime.ErrInvalidArgument, fmt.Sprintf("UDF '%s' expects %d arguments, got %d", udf.Name, udf.Params, argCount))
	}

	if udf.Registers <= 0 {
		return runtime.Error(runtime.ErrInvalidOperation, fmt.Sprintf("UDF '%s' has invalid register window", udf.Name))
	}

	frame := &vm.frames[len(vm.frames)-1]
	frame.fnID = fnID

	if cap(reg) >= udf.Registers {
		reg = reg[:udf.Registers]
		clear(reg)
		if len(args) > 0 && len(reg) > 1 {
			copy(reg[1:], args)
		}
		vm.registers.Values = reg
	} else {
		newRegs := vm.regPool.get(udf.Registers)
		if len(args) > 0 && len(newRegs) > 1 {
			copy(newRegs[1:], args)
		}
		vm.regPool.put(reg)
		vm.registers.Values = newRegs
	}

	vm.pc = udf.Entry

	return nil
}

func (vm *VM) execUdfCall(op bytecode.Opcode, dst, src1, src2 bytecode.Operand) error {
	switch op {
	case bytecode.OpCall, bytecode.OpProtectedCall:
		if err := vm.callUdf(op, dst, src1, src2); err != nil {
			if err := vm.setCallResult(op, dst, runtime.None, err); err != nil {
				if vm.unwindToProtected() {
					return nil
				}

				return err
			}
		}

		return nil
	case bytecode.OpTailCall:
		if err := vm.tailCallUdf(dst, src1, src2); err != nil {
			if vm.unwindToProtected() {
				return nil
			}

			return err
		}

		return nil
	default:
		return runtime.Error(runtime.ErrUnexpected, "invalid udf call opcode")
	}
}

func isProtectedUdfCall(op bytecode.Opcode) bool {
	switch op {
	case bytecode.OpProtectedCall:
		return true
	default:
		return false
	}
}
