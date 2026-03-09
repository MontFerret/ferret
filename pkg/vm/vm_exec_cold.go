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
func (exec *execState) execColdOps(
	ctx context.Context,
	op bytecode.Opcode,
	dst, src1, src2 bytecode.Operand,
) (bool, error) {
	reg := exec.registers.Values
	constants := exec.vm.program.Constants

	switch op {
	case bytecode.OpClose:
		val, ok := reg[dst].(io.Closer)
		reg[dst] = runtime.None

		if ok {
			if err := exec.errors.handle(val.Close()); err != nil {
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

		return true, exec.errors.protected(err)
	case bytecode.OpApplyQuery:
		src := readOperandValue(reg, constants, src1)
		descriptor := readOperandValue(reg, constants, src2)
		out, err := operators.ApplyQuery(ctx, src, descriptor)
		if err := exec.errors.setOrCatch(dst, out, err); err != nil {
			return true, exec.errors.protected(err)
		}

		return true, nil
	case bytecode.OpLike:
		res, err := operators.Like(reg[src1], reg[src2])
		if err == nil {
			reg[dst] = res
			return true, nil
		}

		return true, exec.errors.protected(err)
	case bytecode.OpRegexp:
		r, err := exec.regexpCached(exec.pc-1, reg[src2])
		if err == nil {
			reg[dst] = r.Match(reg[src1])
			return true, nil
		}

		if err := exec.errors.handleWithCatch(err, func() {
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
		return true, exec.errors.setOrCatch(dst, res, err)
	case bytecode.OpStream:
		return true, exec.execStreamOp(ctx, dst, src1, src2)
	case bytecode.OpStreamIter:
		return true, exec.execStreamIterOp(ctx, dst, src1, src2)
	case bytecode.OpDispatch:
		dispatcher, eventName, payload, options, err := exec.castDispatchArgs(ctx, reg[dst], reg[src1], reg[src2])
		if err != nil {
			return true, exec.errors.setOrCatch(dst, runtime.None, err)
		}

		out, err := dispatcher.Dispatch(ctx, runtime.DispatchEvent{
			Name:    eventName,
			Payload: payload,
			Options: options,
		})
		if out == nil {
			out = runtime.None
		}

		return true, exec.errors.setOrCatch(dst, out, err)
	case bytecode.OpFail:
		if !dst.IsConstant() {
			return true, exec.errors.handle(runtime.Error(runtime.ErrInvalidOperation, "FAIL expects a constant string message"))
		}

		idx := dst.Constant()
		if idx < 0 || idx >= len(constants) {
			return true, exec.errors.handle(runtime.Error(runtime.ErrInvalidOperation, "FAIL expects a valid constant string message"))
		}

		msg, ok := constants[idx].(runtime.String)
		if !ok {
			return true, exec.errors.handle(runtime.TypeErrorOf(constants[idx], runtime.TypeString))
		}

		return true, exec.errors.handle(runtime.Error(runtime.ErrInvalidOperation, msg.String()))
	case bytecode.OpSleep:
		dur, err := runtime.ToInt(ctx, reg[dst])
		if err != nil {
			return true, exec.errors.handle(err)
		}

		return true, exec.errors.protected(data.Sleep(ctx, dur))
	default:
		return false, nil
	}
}
