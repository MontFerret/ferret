package vm

import (
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/frame"
)

func (s *execState) resolveUdfID(val runtime.Value) (int, error) {
	idVal, ok := val.(runtime.Int)
	if !ok {
		return -1, ErrInvalidFunctionName
	}

	return int(idVal), nil
}

func (s *execState) udfByID(id int) (*bytecode.UDF, error) {
	if id < 0 || s.program == nil || id >= len(s.program.Functions.UserDefined) {
		return nil, ErrUnresolvedFunction
	}

	return &s.program.Functions.UserDefined[id], nil
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

func (s *execState) callUdf(op bytecode.Opcode, dst, src1, src2 bytecode.Operand) error {
	reg := s.registers.Values

	fnID, err := s.resolveUdfID(reg[dst])
	if err != nil {
		return err
	}

	udf, err := s.udfByID(fnID)
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

	newRegs := s.frames.AcquireRegisters(udf.Registers)
	copyUdfArgsToUdfRegisters(newRegs, reg, argStart, argCount)

	s.frames.Push(frame.CallFrame{
		ReturnPC:   s.pc,
		ReturnDest: dst,
		Registers:  s.registers.Values,
		Protected:  bytecode.IsProtectedUdfCall(op),
		FnID:       fnID,
	})
	s.registers.Values = newRegs
	s.pc = udf.Entry

	return nil
}

func (s *execState) tailCallUdf(dst, src1, src2 bytecode.Operand) error {
	if s.frames.Len() == 0 {
		return ErrUnresolvedFunction
	}

	reg := s.registers.Values
	fnID, err := s.resolveUdfID(reg[dst])
	if err != nil {
		return err
	}

	udf, err := s.udfByID(fnID)
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

	if ok := s.frames.SetTopFnID(fnID); !ok {
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

		s.registers.Values = reg
	} else {
		newRegs := s.frames.AcquireRegisters(udf.Registers)

		if len(args) > 0 && len(newRegs) > 1 {
			copy(newRegs[1:], args)
		}

		s.frames.ReleaseRegisters(reg)
		s.registers.Values = newRegs
	}

	s.pc = udf.Entry

	return nil
}
