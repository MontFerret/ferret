package vm

import (
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/frame"
)

func copyUdfArgsToUdfRegisters(dst, src []runtime.Value, start, count int) {
	if len(dst) <= 1 || len(src) == 0 || count <= 0 {
		return
	}

	if maxCount := len(dst) - 1; count > maxCount {
		count = maxCount
	}

	if count <= 0 {
		return
	}

	copy(dst[1:1+count], src[start:start+count])
}

func collectUdfArgsInto(dst, src []runtime.Value, start, count int) int {
	if count <= 0 {
		return 0
	}

	if count > len(dst) {
		count = len(dst)
	}

	copy(dst[:count], src[start:start+count])

	return count
}

func callUdf(s *execState, desc *callDescriptor, udf *bytecode.UDF) error {
	reg := s.registers

	if desc.ArgCount != udf.Params {
		return runtime.Error(runtime.ErrInvalidArgument, fmt.Sprintf("UDF '%s' expects %d arguments, got %d", udf.Name, udf.Params, desc.ArgCount))
	}

	if udf.Registers <= 0 {
		return runtime.Error(runtime.ErrInvalidOperation, fmt.Sprintf("UDF '%s' has invalid register window", desc.DisplayName))
	}

	newRegs := s.frames.AcquireRegisters(udf.Registers)
	copyUdfArgsToUdfRegisters(newRegs, reg, desc.ArgStart, desc.ArgCount)

	s.frames.Push(frame.CallFrame{
		ReturnPC:         s.pc,
		ReturnDest:       desc.Dst,
		Registers:        s.registers,
		RecoveryBoundary: desc.RecoveryBoundary,
		FnID:             desc.ID,
		FnName:           desc.DisplayName,
		CallSitePC:       s.pc - 1,
		HasCallSite:      true,
	})
	s.registers = newRegs
	s.pc = udf.Entry

	return nil
}

func tailCallUdf(s *execState, desc *callDescriptor, udf *bytecode.UDF) error {
	reg := s.registers
	argStart := desc.ArgStart
	argCount := desc.ArgCount

	if udf.Params != desc.ArgCount {
		return runtime.Error(runtime.ErrInvalidArgument, fmt.Sprintf("UDF '%s' expects %d arguments, got %d", desc.DisplayName, udf.Params, desc.ArgCount))
	}

	if udf.Registers <= 0 {
		return runtime.Error(runtime.ErrInvalidOperation, fmt.Sprintf("UDF '%s' has invalid register window", desc.DisplayName))
	}

	if ok := s.frames.SetTopCall(desc.ID, desc.DisplayName, desc.CallSitePC); !ok {
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
		for i := range reg {
			reg[i] = runtime.None
		}

		if len(args) > 0 && len(reg) > 1 {
			copy(reg[1:], args)
		}

		s.registers = reg
	} else {
		newRegs := s.frames.AcquireRegisters(udf.Registers)

		if len(args) > 0 && len(newRegs) > 1 {
			copy(newRegs[1:], args)
		}

		s.frames.ReleaseRegisters(reg)
		s.registers = newRegs
	}

	s.pc = udf.Entry

	return nil
}
