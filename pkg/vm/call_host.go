package vm

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/mem"
)

func callCachedHostFunction(
	ctx context.Context,
	desc *callDescriptor,
	cacheFn *mem.CachedHostFunction,
	reg []runtime.Value,
	scratch *mem.Scratch,
) (runtime.Value, error) {
	if cacheFn == nil || !cacheFn.Bound {
		return nil, ErrUnresolvedFunction
	}

	start := desc.ArgStart
	count := desc.ArgCount
	fnv := cacheFn.FnV

	switch count {
	case 0:
		if cacheFn.Fn0 != nil {
			return cacheFn.Fn0(ctx)
		}

		if fnv != nil {
			return fnv(ctx)
		}
	case 1:
		arg0 := reg[start]

		if cacheFn.Fn1 != nil {
			return cacheFn.Fn1(ctx, arg0)
		}

		if fnv != nil {
			return fnv(ctx, arg0)
		}
	case 2:
		arg0 := reg[start]
		arg1 := reg[start+1]

		if cacheFn.Fn2 != nil {
			return cacheFn.Fn2(ctx, arg0, arg1)
		}

		if fnv != nil {
			return fnv(ctx, arg0, arg1)
		}
	case 3:
		arg0 := reg[start]
		arg1 := reg[start+1]
		arg2 := reg[start+2]

		if cacheFn.Fn3 != nil {
			return cacheFn.Fn3(ctx, arg0, arg1, arg2)
		}

		if fnv != nil {
			return fnv(ctx, arg0, arg1, arg2)
		}
	case 4:
		arg0 := reg[start]
		arg1 := reg[start+1]
		arg2 := reg[start+2]
		arg3 := reg[start+3]

		if cacheFn.Fn4 != nil {
			return cacheFn.Fn4(ctx, arg0, arg1, arg2, arg3)
		}

		if fnv != nil {
			return fnv(ctx, arg0, arg1, arg2, arg3)
		}
	default:
		if fnv != nil {
			args := stageHostCallArgs(scratch, reg, start, count)
			return fnv(ctx, args...)
		}
	}

	return nil, ErrUnresolvedFunction
}

func stageHostCallArgs(scratch *mem.Scratch, reg []runtime.Value, start, count int) []runtime.Value {
	if count <= 0 {
		return nil
	}

	scratch.ResizeHostArgs(count)
	args := scratch.HostArgs[:count]
	copy(args, reg[start:start+count])

	return args
}
