package vm

import (
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/frame"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/mem"
)

func copyUdfArgsToUdfRegisters(dst, src []runtime.Value, start, count int) {
	count = udfArgCountForRegisters(len(dst), count)
	if count <= 0 {
		return
	}

	copy(dst[1:1+count], src[start:start+count])
}

func udfArgCountForRegisters(registers, count int) int {
	if registers <= 1 || count <= 0 {
		return 0
	}

	if maxCount := registers - 1; count > maxCount {
		count = maxCount
	}

	return count
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
	if desc.ArgCount != udf.Params {
		return runtime.Error(runtime.ErrInvalidArgument, fmt.Sprintf("UDF '%s' expects %d arguments, got %d", desc.DisplayName, udf.Params, desc.ArgCount))
	}

	if udf.Registers <= 0 {
		return runtime.Error(runtime.ErrInvalidOperation, fmt.Sprintf("UDF '%s' has invalid register window", desc.DisplayName))
	}

	s.enterUdfCall(desc, udf)

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

	copyCount := udfArgCountForRegisters(udf.Registers, argCount)
	if copyCount > 0 {
		if copyCount <= len(stackArgs) {
			count := collectUdfArgsInto(stackArgs[:copyCount], reg, argStart, copyCount)
			args = stackArgs[:count]
		} else {
			heapArgs = make([]runtime.Value, copyCount)
			count := collectUdfArgsInto(heapArgs, reg, argStart, copyCount)
			args = heapArgs[:count]
		}
	}

	if cap(reg) >= udf.Registers {
		var ownedArgs mem.OwnedResources
		s.owned.ExtractMany(args, &ownedArgs)
		s.owned.DrainTo(&s.deferred)

		reg = reg[:udf.Registers]
		for i := range reg {
			reg[i] = runtime.None
		}

		if len(args) > 0 && len(reg) > 1 {
			copy(reg[1:], args)
		}

		s.registers = reg
		s.owned = ownedArgs
	} else {
		newRegs := s.windows.Acquire(udf.Registers)
		copyUdfArgsToUdfRegisters(newRegs, reg, argStart, copyCount)

		var ownedArgs mem.OwnedResources
		s.owned.ExtractMany(args, &ownedArgs)
		s.owned.DrainTo(&s.deferred)

		s.windows.Release(reg)
		s.registers = newRegs
		s.owned = ownedArgs
	}

	// Rebuild alias counts from surviving owned args.
	// Must iterate args (not owned) because multiple args may alias
	// the same closer, and each register slot needs its own count.
	s.aliases.Reset()

	for _, arg := range args {
		key, _, ok := mem.ResourceKeyOf(arg)
		if ok && s.owned.Owns(arg) {
			s.aliases.Inc(key)
		}
	}

	s.pc = udf.Entry

	return nil
}

func (s *execState) enterUdfCall(desc *callDescriptor, udf *bytecode.UDF) {
	newRegs := s.windows.Acquire(udf.Registers)
	copyUdfArgsToUdfRegisters(newRegs, s.registers, desc.ArgStart, desc.ArgCount)

	s.frames.Push(frame.CallFrame{
		ReturnPC:         s.pc,
		ReturnDest:       desc.Dst,
		CallerRegisters:  s.registers,
		OwnedResources:   s.owned,
		Aliases:          s.aliases,
		RecoveryBoundary: desc.RecoveryBoundary,
		FnID:             desc.ID,
		FnName:           desc.DisplayName,
		CallSitePC:       desc.CallSitePC,
		HasCallSite:      true,
	})
	s.owned = mem.OwnedResources{}
	s.aliases = mem.AliasTracker{}
	s.registers = newRegs
	s.pc = udf.Entry
}
