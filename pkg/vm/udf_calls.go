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

func (vm *VM) udfArgCount(op bytecode.Opcode, src1, src2 bytecode.Operand) int {
	switch op {
	case bytecode.OpCall0, bytecode.OpProtectedCall0, bytecode.OpTailCall0:
		return 0
	case bytecode.OpCall1, bytecode.OpProtectedCall1, bytecode.OpTailCall1:
		return 1
	case bytecode.OpCall2, bytecode.OpProtectedCall2, bytecode.OpTailCall2:
		return 2
	case bytecode.OpCall3, bytecode.OpProtectedCall3, bytecode.OpTailCall3:
		return 3
	case bytecode.OpCall4, bytecode.OpProtectedCall4, bytecode.OpTailCall4:
		return 4
	default:
		if !src1.IsRegister() || !src2.IsRegister() {
			return 0
		}

		return src2.Register() - src1.Register() + 1
	}
}

func (vm *VM) copyUdfArgs(op bytecode.Opcode, dst []runtime.Value, src []runtime.Value, src1, src2 bytecode.Operand) {
	switch op {
	case bytecode.OpCall1, bytecode.OpProtectedCall1, bytecode.OpTailCall1:
		dst[1] = src[src1]
	case bytecode.OpCall2, bytecode.OpProtectedCall2, bytecode.OpTailCall2:
		dst[1] = src[src1]
		dst[2] = src[src2]
	case bytecode.OpCall3, bytecode.OpProtectedCall3, bytecode.OpTailCall3:
		start := src1.Register()
		dst[1] = src[start]
		dst[2] = src[start+1]
		dst[3] = src[start+2]
	case bytecode.OpCall4, bytecode.OpProtectedCall4, bytecode.OpTailCall4:
		start := src1.Register()
		dst[1] = src[start]
		dst[2] = src[start+1]
		dst[3] = src[start+2]
		dst[4] = src[start+3]
	case bytecode.OpCall, bytecode.OpProtectedCall, bytecode.OpTailCall:
		start := src1.Register()
		end := src2.Register()
		if start <= 0 || end < start {
			return
		}
		for i := 0; i <= end-start && i+1 < len(dst); i++ {
			dst[1+i] = src[start+i]
		}
	default:
		// OpCall0, OpProtectedCall0, OpTailCall0 have no arguments.
	}
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

	argCount := vm.udfArgCount(op, src1, src2)
	if udf.Params != argCount {
		return runtime.Error(runtime.ErrInvalidArgument, fmt.Sprintf("UDF '%s' expects %d arguments, got %d", udf.Name, udf.Params, argCount))
	}

	if udf.Registers <= 0 {
		return runtime.Error(runtime.ErrInvalidOperation, fmt.Sprintf("UDF '%s' has invalid register window", udf.Name))
	}

	newRegs := make([]runtime.Value, udf.Registers)
	vm.copyUdfArgs(op, newRegs, reg, src1, src2)

	vm.pushFrame(vm.pc, dst, isProtectedUdfCall(op), fnID)
	vm.registers.Values = newRegs
	vm.pc = udf.Entry

	return nil
}

func (vm *VM) tailCallUdf(op bytecode.Opcode, dst, src1, src2 bytecode.Operand) error {
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

	argCount := vm.udfArgCount(op, src1, src2)
	if udf.Params != argCount {
		return runtime.Error(runtime.ErrInvalidArgument, fmt.Sprintf("UDF '%s' expects %d arguments, got %d", udf.Name, udf.Params, argCount))
	}

	if udf.Registers <= 0 {
		return runtime.Error(runtime.ErrInvalidOperation, fmt.Sprintf("UDF '%s' has invalid register window", udf.Name))
	}

	newRegs := make([]runtime.Value, udf.Registers)
	vm.copyUdfArgs(op, newRegs, reg, src1, src2)

	frame := &vm.frames[len(vm.frames)-1]
	frame.fnID = fnID
	vm.registers.Values = newRegs
	vm.pc = udf.Entry

	return nil
}

func isProtectedUdfCall(op bytecode.Opcode) bool {
	switch op {
	case bytecode.OpProtectedCall, bytecode.OpProtectedCall0, bytecode.OpProtectedCall1, bytecode.OpProtectedCall2, bytecode.OpProtectedCall3, bytecode.OpProtectedCall4:
		return true
	default:
		return false
	}
}
