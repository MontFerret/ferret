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

func callCachedHostFunction(
	ctx context.Context,
	cacheFn *mem.CachedHostFunction,
	reg []runtime.Value,
	target runtime.Value,
	src1, src2 bytecode.Operand,
) (runtime.Value, error) {
	if cacheFn == nil {
		if _, ok := target.(runtime.String); !ok {
			return nil, ErrInvalidFunctionName
		}

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
			argCount := end - start + 1
			var stackArgs [8]runtime.Value
			if argCount <= len(stackArgs) {
				copy(stackArgs[:argCount], reg[start:start+argCount])

				switch argCount {
				case 5:
					return cacheFn.FnV(ctx, stackArgs[0], stackArgs[1], stackArgs[2], stackArgs[3], stackArgs[4])
				case 6:
					return cacheFn.FnV(ctx, stackArgs[0], stackArgs[1], stackArgs[2], stackArgs[3], stackArgs[4], stackArgs[5])
				case 7:
					return cacheFn.FnV(ctx, stackArgs[0], stackArgs[1], stackArgs[2], stackArgs[3], stackArgs[4], stackArgs[5], stackArgs[6])
				case 8:
					return cacheFn.FnV(ctx, stackArgs[0], stackArgs[1], stackArgs[2], stackArgs[3], stackArgs[4], stackArgs[5], stackArgs[6], stackArgs[7])
				}
			}

			heapArgs := make([]runtime.Value, argCount)
			copy(heapArgs, reg[start:start+argCount])
			return cacheFn.FnV(ctx, heapArgs...)
		}
	}

	return nil, ErrUnresolvedFunction
}
