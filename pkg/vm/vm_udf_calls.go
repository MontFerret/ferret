package vm

import (
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/frame"
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

func udfArgInfo(src1, src2 bytecode.Operand) (int, int) {
	start, end, ok := udfArgRange(src1, src2)
	if !ok {
		return 0, 0
	}

	return start, end - start + 1
}

func clampUdfArgCount(srcLen, start, count int) int {
	if count <= 0 || srcLen == 0 || start >= srcLen {
		return 0
	}

	if maxCount := srcLen - start; count > maxCount {
		count = maxCount
	}

	return count
}

func copyUdfArgsToUdfRegisters(dst, src []runtime.Value, start, count int) {
	if len(dst) <= 1 || len(src) == 0 || count <= 0 {
		return
	}

	if maxCount := len(dst) - 1; count > maxCount {
		count = maxCount
	}

	count = clampUdfArgCount(len(src), start, count)
	if count <= 0 {
		return
	}

	copy(dst[1:1+count], src[start:start+count])
}

func collectUdfArgsInto(dst, src []runtime.Value, start, count int) int {
	count = clampUdfArgCount(len(src), start, count)
	if count <= 0 {
		return 0
	}

	if count > len(dst) {
		count = len(dst)
	}

	copy(dst[:count], src[start:start+count])

	return count
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

	argStart, argCount := udfArgInfo(src1, src2)
	if udf.Params != argCount {
		return runtime.Error(runtime.ErrInvalidArgument, fmt.Sprintf("UDF '%s' expects %d arguments, got %d", udf.Name, udf.Params, argCount))
	}

	if udf.Registers <= 0 {
		return runtime.Error(runtime.ErrInvalidOperation, fmt.Sprintf("UDF '%s' has invalid register window", udf.Name))
	}

	newRegs := vm.frames.GetRegisters(udf.Registers)
	copyUdfArgsToUdfRegisters(newRegs, reg, argStart, argCount)

	vm.frames.Push(frame.CallFrame{
		ReturnPC:   vm.pc,
		ReturnDest: dst,
		Registers:  vm.registers.Values,
		Protected:  isProtectedUdfCall(op),
		FnID:       fnID,
	})
	vm.registers.Values = newRegs
	vm.pc = udf.Entry

	return nil
}

func (vm *VM) tailCallUdf(dst, src1, src2 bytecode.Operand) error {
	if vm.frames.Len() == 0 {
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

	argStart, argCount := udfArgInfo(src1, src2)
	if udf.Params != argCount {
		return runtime.Error(runtime.ErrInvalidArgument, fmt.Sprintf("UDF '%s' expects %d arguments, got %d", udf.Name, udf.Params, argCount))
	}

	if udf.Registers <= 0 {
		return runtime.Error(runtime.ErrInvalidOperation, fmt.Sprintf("UDF '%s' has invalid register window", udf.Name))
	}

	currentFrame := vm.frames.Top()
	if currentFrame == nil {
		return ErrUnresolvedFunction
	}
	currentFrame.FnID = fnID

	var (
		args      []runtime.Value
		heapArgs  []runtime.Value
		stackArgs [8]runtime.Value
	)

	if argCount > 0 {
		if argCount <= len(stackArgs) {
			count := collectUdfArgsInto(stackArgs[:argCount], reg, argStart, argCount)
			args = stackArgs[:count]
		} else {
			heapArgs = make([]runtime.Value, argCount)
			count := collectUdfArgsInto(heapArgs, reg, argStart, argCount)
			args = heapArgs[:count]
		}
	}

	if cap(reg) >= udf.Registers {
		reg = reg[:udf.Registers]
		clear(reg)
		if len(args) > 0 && len(reg) > 1 {
			copy(reg[1:], args)
		}
		vm.registers.Values = reg
	} else {
		newRegs := vm.frames.GetRegisters(udf.Registers)
		if len(args) > 0 && len(newRegs) > 1 {
			copy(newRegs[1:], args)
		}
		vm.frames.PutRegisters(reg)
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
