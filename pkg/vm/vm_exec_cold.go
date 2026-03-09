package vm

import (
	"context"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/operators"
)

// execColdOps keeps low-frequency and error-heavy instructions out of the main dispatch hot path.
func (vm *VM) execColdOps(
	ctx context.Context,
	op bytecode.Opcode,
	dst, src1, src2 bytecode.Operand,
	reg []runtime.Value,
	constants []runtime.Value,
) (bool, error) {
	switch op {
	case bytecode.OpClose:
		val, ok := reg[dst].(io.Closer)
		reg[dst] = runtime.None

		if ok {
			if err := vm.errors.handle(val.Close()); err != nil {
				return true, err
			}
		}

		return true, nil
	case bytecode.OpLoadRange:
		res, err := operators.ToRange(ctx, reg[src1], reg[src2])
		if err == nil {
			reg[dst] = res
			return true, nil
		}

		return true, vm.errors.protected(err)
	case bytecode.OpApplyQuery:
		src := readOperandValue(reg, constants, src1)
		descriptor := readOperandValue(reg, constants, src2)
		out, err := operators.ApplyQuery(ctx, src, descriptor)
		if err := vm.errors.setOrCatch(dst, out, err); err != nil {
			return true, vm.errors.protected(err)
		}

		return true, nil
	case bytecode.OpLike:
		res, err := operators.Like(reg[src1], reg[src2])
		if err == nil {
			reg[dst] = res
			return true, nil
		}

		return true, vm.errors.protected(err)
	case bytecode.OpRegexp:
		r, err := vm.regexpCached(vm.pc-1, reg[src2])
		if err == nil {
			reg[dst] = r.Match(reg[src1])
			return true, nil
		}

		if err := vm.errors.handleWithCatch(err, func() {
			reg[dst] = runtime.False
		}); err != nil {
			return true, err
		}

		return true, nil
	case bytecode.OpFlatten:
		depth := src2.Register()
		if depth < 1 {
			depth = 1
		}

		res, err := operators.Flatten(ctx, reg[src1], depth)
		return true, vm.errors.setOrCatch(dst, res, err)
	case bytecode.OpStream:
		return true, vm.execStreamOp(ctx, dst, src1, src2, reg)
	case bytecode.OpStreamIter:
		return true, vm.execStreamIterOp(ctx, dst, src1, src2, reg)
	case bytecode.OpDispatch:
		dispatcher, eventName, payload, options, err := vm.castDispatchArgs(ctx, reg[dst], reg[src1], reg[src2])
		if err != nil {
			return true, vm.errors.setOrCatch(dst, runtime.None, err)
		}

		out, err := dispatcher.Dispatch(ctx, runtime.DispatchEvent{
			Name:    eventName,
			Payload: payload,
			Options: options,
		})
		if out == nil {
			out = runtime.None
		}

		return true, vm.errors.setOrCatch(dst, out, err)
	case bytecode.OpFail:
		if !dst.IsConstant() {
			return true, vm.errors.handle(runtime.Error(runtime.ErrInvalidOperation, "FAIL expects a constant string message"))
		}

		idx := dst.Constant()
		if idx < 0 || idx >= len(constants) {
			return true, vm.errors.handle(runtime.Error(runtime.ErrInvalidOperation, "FAIL expects a valid constant string message"))
		}

		msg, ok := constants[idx].(runtime.String)
		if !ok {
			return true, vm.errors.handle(runtime.TypeErrorOf(constants[idx], runtime.TypeString))
		}

		return true, vm.errors.handle(runtime.Error(runtime.ErrInvalidOperation, msg.String()))
	case bytecode.OpSleep:
		dur, err := runtime.ToInt(ctx, reg[dst])
		if err != nil {
			return true, vm.errors.handle(err)
		}

		return true, vm.errors.protected(data.Sleep(ctx, dur))
	default:
		return false, nil
	}
}
