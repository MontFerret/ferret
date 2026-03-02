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
		idxs := [2]int{src1.Register()}
		copyUdfArgsInternal(dst, src, 1, 0, 1, false, false, idxs[:1])
	case bytecode.OpCall2, bytecode.OpProtectedCall2, bytecode.OpTailCall2:
		idxs := [2]int{src1.Register(), src2.Register()}
		copyUdfArgsInternal(dst, src, 1, 0, 2, false, false, idxs[:2])
	case bytecode.OpCall3, bytecode.OpProtectedCall3, bytecode.OpTailCall3:
		start := src1.Register()
		copyUdfArgsInternal(dst, src, 1, start, 3, false, false, nil)
	case bytecode.OpCall4, bytecode.OpProtectedCall4, bytecode.OpTailCall4:
		start := src1.Register()
		copyUdfArgsInternal(dst, src, 1, start, 4, false, false, nil)
	case bytecode.OpCall, bytecode.OpProtectedCall, bytecode.OpTailCall:
		start := src1.Register()
		end := src2.Register()
		if start <= 0 || end < start {
			return
		}
		copyUdfArgsInternal(dst, src, 1, start, end-start+1, true, false, nil)
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

	newRegs := vm.regPool.get(udf.Registers)
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

	frame := &vm.frames[len(vm.frames)-1]
	frame.fnID = fnID

	if cap(reg) >= udf.Registers {
		oldLen := len(reg)
		vm.copyUdfArgsInPlace(op, reg, src1, src2)

		if oldLen > udf.Registers {
			for i := udf.Registers; i < oldLen; i++ {
				reg[i] = nil
			}
		}

		if oldLen != udf.Registers {
			reg = reg[:udf.Registers]
		}

		clearUdfRegsExceptArgs(reg, argCount)
		vm.registers.Values = reg
	} else {
		newRegs := vm.regPool.get(udf.Registers)
		vm.copyUdfArgs(op, newRegs, reg, src1, src2)
		vm.regPool.put(reg)
		vm.registers.Values = newRegs
	}
	vm.pc = udf.Entry

	return nil
}

func (vm *VM) copyUdfArgsInPlace(op bytecode.Opcode, reg []runtime.Value, src1, src2 bytecode.Operand) {
	switch op {
	case bytecode.OpCall1, bytecode.OpProtectedCall1, bytecode.OpTailCall1:
		idxs := [2]int{src1.Register()}
		copyUdfArgsInternal(reg, reg, 1, 0, 1, false, true, idxs[:1])
	case bytecode.OpCall2, bytecode.OpProtectedCall2, bytecode.OpTailCall2:
		idxs := [2]int{src1.Register(), src2.Register()}
		copyUdfArgsInternal(reg, reg, 1, 0, 2, false, true, idxs[:2])
	case bytecode.OpCall3, bytecode.OpProtectedCall3, bytecode.OpTailCall3:
		start := src1.Register()
		copyUdfArgsInternal(reg, reg, 1, start, 3, false, true, nil)
	case bytecode.OpCall4, bytecode.OpProtectedCall4, bytecode.OpTailCall4:
		start := src1.Register()
		copyUdfArgsInternal(reg, reg, 1, start, 4, false, true, nil)
	case bytecode.OpCall, bytecode.OpProtectedCall, bytecode.OpTailCall:
		start := src1.Register()
		end := src2.Register()
		if start <= 0 || end < start {
			return
		}
		copyUdfArgsInternal(reg, reg, 1, start, end-start+1, true, true, nil)
	default:
		// OpCall0, OpProtectedCall0, OpTailCall0 have no arguments.
	}
}

func copyUdfArgsInternal(dst, src []runtime.Value, dstStart, srcStart, count int, bounded bool, sameSlice bool, indices []int) {
	if count <= 0 {
		return
	}

	if indices != nil {
		if sameSlice {
			var tmp [4]runtime.Value
			if count <= len(tmp) {
				for i := 0; i < count; i++ {
					tmp[i] = src[indices[i]]
				}
				for i := 0; i < count; i++ {
					if bounded && dstStart+i >= len(dst) {
						continue
					}
					dst[dstStart+i] = tmp[i]
				}
				return
			}

			values := make([]runtime.Value, count)
			for i := 0; i < count; i++ {
				values[i] = src[indices[i]]
			}
			for i := 0; i < count; i++ {
				if bounded && dstStart+i >= len(dst) {
					continue
				}
				dst[dstStart+i] = values[i]
			}
			return
		}

		for i := 0; i < count; i++ {
			if bounded && dstStart+i >= len(dst) {
				return
			}
			dst[dstStart+i] = src[indices[i]]
		}
		return
	}

	dstEnd := dstStart + count - 1
	srcEnd := srcStart + count - 1

	if sameSlice && srcStart <= dstEnd && dstStart <= srcEnd {
		for i := count - 1; i >= 0; i-- {
			if bounded && dstStart+i >= len(dst) {
				continue
			}

			dst[dstStart+i] = src[srcStart+i]
		}

		return
	}

	if bounded {
		for i := 0; i < count && dstStart+i < len(dst); i++ {
			dst[dstStart+i] = src[srcStart+i]
		}

		return
	}

	for i := 0; i < count; i++ {
		dst[dstStart+i] = src[srcStart+i]
	}
}

func (vm *VM) execUdfCall(op bytecode.Opcode, dst, src1, src2 bytecode.Operand) error {
	switch op {
	case bytecode.OpCall, bytecode.OpProtectedCall,
		bytecode.OpCall0, bytecode.OpProtectedCall0,
		bytecode.OpCall1, bytecode.OpProtectedCall1,
		bytecode.OpCall2, bytecode.OpProtectedCall2,
		bytecode.OpCall3, bytecode.OpProtectedCall3,
		bytecode.OpCall4, bytecode.OpProtectedCall4:
		if op == bytecode.OpCall0 || op == bytecode.OpProtectedCall0 {
			src1, src2 = 0, 0
		} else if op == bytecode.OpCall3 || op == bytecode.OpProtectedCall3 || op == bytecode.OpCall4 || op == bytecode.OpProtectedCall4 {
			src2 = 0
		}

		if err := vm.callUdf(op, dst, src1, src2); err != nil {
			if err := vm.setCallResult(op, dst, runtime.None, err); err != nil {
				if vm.unwindToProtected() {
					return nil
				}

				return err
			}
		}

		return nil
	case bytecode.OpTailCall, bytecode.OpTailCall0, bytecode.OpTailCall1, bytecode.OpTailCall2, bytecode.OpTailCall3, bytecode.OpTailCall4:
		if op == bytecode.OpTailCall0 {
			src1, src2 = 0, 0
		} else if op == bytecode.OpTailCall3 || op == bytecode.OpTailCall4 {
			src2 = 0
		}

		if err := vm.tailCallUdf(op, dst, src1, src2); err != nil {
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

func clearUdfRegsExceptArgs(reg []runtime.Value, argCount int) {
	if len(reg) == 0 {
		return
	}

	if argCount < 0 {
		argCount = 0
	}

	reg[0] = nil

	if argCount >= len(reg)-1 {
		return
	}

	for i := argCount + 1; i < len(reg); i++ {
		reg[i] = nil
	}
}

func isProtectedUdfCall(op bytecode.Opcode) bool {
	switch op {
	case bytecode.OpProtectedCall, bytecode.OpProtectedCall0, bytecode.OpProtectedCall1, bytecode.OpProtectedCall2, bytecode.OpProtectedCall3, bytecode.OpProtectedCall4:
		return true
	default:
		return false
	}
}
