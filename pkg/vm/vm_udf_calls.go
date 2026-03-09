package vm

import (
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/frame"
)

func (exec *execState) resolveUdfID(val runtime.Value) (int, error) {
	idVal, ok := val.(runtime.Int)
	if !ok {
		return -1, ErrInvalidFunctionName
	}

	return int(idVal), nil
}

func (exec *execState) udfByID(id int) (*bytecode.UDF, error) {
	if id < 0 || exec.vm.program == nil || id >= len(exec.vm.program.Functions.UserDefined) {
		return nil, ErrUnresolvedFunction
	}

	return &exec.vm.program.Functions.UserDefined[id], nil
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

func (exec *execState) callUdf(op bytecode.Opcode, dst, src1, src2 bytecode.Operand) error {
	reg := exec.registers.Values

	fnID, err := exec.resolveUdfID(reg[dst])
	if err != nil {
		return err
	}

	udf, err := exec.udfByID(fnID)
	if err != nil {
		return err
	}

	argStart, argCount := callArgInfo(src1, src2)
	if udf.Params != argCount {
		return runtime.Error(runtime.ErrInvalidArgument, fmt.Sprintf("UDF '%s' expects %d arguments, got %d", udf.Name, udf.Params, argCount))
	}

	if udf.Registers <= 0 {
		return runtime.Error(runtime.ErrInvalidOperation, fmt.Sprintf("UDF '%s' has invalid register window", udf.Name))
	}

	newRegs := exec.vm.frames.AcquireRegisters(udf.Registers)
	copyUdfArgsToUdfRegisters(newRegs, reg, argStart, argCount)

	exec.vm.frames.Push(frame.CallFrame{
		ReturnPC:   exec.pc,
		ReturnDest: dst,
		Registers:  exec.registers.Values,
		Protected:  bytecode.IsProtectedCallOpcode(op),
		FnID:       fnID,
	})
	exec.registers.Values = newRegs
	exec.pc = udf.Entry

	return nil
}

func (exec *execState) tailCallUdf(dst, src1, src2 bytecode.Operand) error {
	if exec.vm.frames.Len() == 0 {
		return ErrUnresolvedFunction
	}

	reg := exec.registers.Values
	fnID, err := exec.resolveUdfID(reg[dst])
	if err != nil {
		return err
	}

	udf, err := exec.udfByID(fnID)
	if err != nil {
		return err
	}

	argStart, argCount := callArgInfo(src1, src2)
	if udf.Params != argCount {
		return runtime.Error(runtime.ErrInvalidArgument, fmt.Sprintf("UDF '%s' expects %d arguments, got %d", udf.Name, udf.Params, argCount))
	}

	if udf.Registers <= 0 {
		return runtime.Error(runtime.ErrInvalidOperation, fmt.Sprintf("UDF '%s' has invalid register window", udf.Name))
	}

	if ok := exec.vm.frames.SetTopFnID(fnID); !ok {
		return ErrUnresolvedFunction
	}

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
		exec.registers.Values = reg
	} else {
		newRegs := exec.vm.frames.AcquireRegisters(udf.Registers)
		if len(args) > 0 && len(newRegs) > 1 {
			copy(newRegs[1:], args)
		}
		exec.vm.frames.ReleaseRegisters(reg)
		exec.registers.Values = newRegs
	}

	exec.pc = udf.Entry

	return nil
}

func (exec *execState) execUdfCall(op bytecode.Opcode, dst, src1, src2 bytecode.Operand) error {
	switch op {
	case bytecode.OpCall, bytecode.OpProtectedCall:
		if err := exec.callUdf(op, dst, src1, src2); err != nil {
			if err := exec.errors.setCallResult(op, dst, runtime.None, err); err != nil {
				if exec.unwindToProtected() {
					return nil
				}

				return err
			}
		}

		return nil
	case bytecode.OpTailCall:
		if err := exec.tailCallUdf(dst, src1, src2); err != nil {
			if exec.unwindToProtected() {
				return nil
			}

			return err
		}

		return nil
	default:
		return runtime.Error(runtime.ErrUnexpected, "invalid udf call opcode")
	}
}
