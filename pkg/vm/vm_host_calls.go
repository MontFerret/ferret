package vm

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/mem"
)

func hostCallArgRange(src1, src2 bytecode.Operand) (int, int, bool) {
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

func hostCallArgs(reg []runtime.Value, start, end int) []runtime.Value {
	size := end - start + 1
	args := make([]runtime.Value, size)

	for i := 0; i < size; i++ {
		args[i] = reg[start+i]
	}

	return args
}

func callCachedHostFunction(
	ctx context.Context,
	cacheFn *mem.CachedHostFunction,
	reg []runtime.Value,
	src1, src2 bytecode.Operand,
) (runtime.Value, error) {
	if cacheFn == nil {
		return nil, ErrUnresolvedFunction
	}

	start, end, hasRange := hostCallArgRange(src1, src2)
	if !hasRange {
		if cacheFn.Fn0 != nil {
			return cacheFn.Fn0(ctx)
		}

		if cacheFn.FnV != nil {
			return cacheFn.FnV(ctx)
		}

		return nil, ErrUnresolvedFunction
	}

	if start < 0 || end >= len(reg) {
		return nil, runtime.Error(runtime.ErrUnexpected, "invalid host call argument range")
	}

	switch end - start + 1 {
	case 1:
		arg0 := reg[start]

		if cacheFn.Fn1 != nil {
			return cacheFn.Fn1(ctx, arg0)
		}

		if cacheFn.FnV != nil {
			return cacheFn.FnV(ctx, arg0)
		}
	case 2:
		arg0 := reg[start]
		arg1 := reg[start+1]

		if cacheFn.Fn2 != nil {
			return cacheFn.Fn2(ctx, arg0, arg1)
		}

		if cacheFn.FnV != nil {
			return cacheFn.FnV(ctx, arg0, arg1)
		}
	case 3:
		arg0 := reg[start]
		arg1 := reg[start+1]
		arg2 := reg[start+2]

		if cacheFn.Fn3 != nil {
			return cacheFn.Fn3(ctx, arg0, arg1, arg2)
		}

		if cacheFn.FnV != nil {
			return cacheFn.FnV(ctx, arg0, arg1, arg2)
		}
	case 4:
		arg0 := reg[start]
		arg1 := reg[start+1]
		arg2 := reg[start+2]
		arg3 := reg[start+3]

		if cacheFn.Fn4 != nil {
			return cacheFn.Fn4(ctx, arg0, arg1, arg2, arg3)
		}

		if cacheFn.FnV != nil {
			return cacheFn.FnV(ctx, arg0, arg1, arg2, arg3)
		}
	default:
		if cacheFn.FnV != nil {
			args := hostCallArgs(reg, start, end)
			return cacheFn.FnV(ctx, args...)
		}
	}

	return nil, ErrUnresolvedFunction
}

func (s *execState) setCallResult(op bytecode.Opcode, dst bytecode.Operand, out runtime.Value, err error) errAction {
	reg := s.registers.Values

	if err == nil {
		reg[dst] = out

		return errOK
	}

	if bytecode.IsProtectedCall(op) {
		reg[dst] = runtime.None

		return errContinue
	}

	return s.setCallResultSlow(dst, err)
}

func (s *execState) setCallResultSlow(dst bytecode.Operand, err error) errAction {
	reg := s.registers.Values

	if catch, ok := s.tryCatch(s.pc); ok {
		reg[dst] = runtime.None

		if catch[2] >= 0 {
			s.pc = catch[2]
		}

		return errContinue
	}

	if s.unwindToProtected() {
		return errContinue
	}

	return errReturn
}
