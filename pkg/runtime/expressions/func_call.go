package expressions

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type FunctionCallExpression struct {
	src          core.SourceMap
	fun          core.Function
	ignoreErrors bool
	args         []core.Expression
}

func NewFunctionCallExpression(
	src core.SourceMap,
	fun core.Function,
	ignoreErrors bool,
	args ...core.Expression,
) (*FunctionCallExpression, error) {
	if fun == nil {
		return nil, core.Error(core.ErrMissedArgument, "function")
	}

	return &FunctionCallExpression{src, fun, ignoreErrors, args}, nil
}

func (e *FunctionCallExpression) Arguments() []core.Expression {
	return e.args
}

func (e *FunctionCallExpression) Function() core.Function {
	if e.ignoreErrors {
		return func(ctx context.Context, args ...core.Value) (core.Value, error) {
			out, err := e.fun(ctx, args...)

			if err != nil {
				return values.None, nil
			}

			return out, nil
		}
	}

	return e.fun
}

func (e *FunctionCallExpression) Exec(ctx context.Context, scope *core.Scope) (core.Value, error) {
	select {
	case <-ctx.Done():
		return values.None, core.ErrTerminated
	default:
		var out core.Value
		var err error

		if len(e.args) == 0 {
			out, err = e.fun(ctx)
		} else {
			args := make([]core.Value, len(e.args))

			for idx, arg := range e.args {
				out, err := arg.Exec(ctx, scope)

				if err != nil {
					return values.None, e.maybeError(core.SourceError(e.src, err))
				}

				args[idx] = out
			}

			out, err = e.fun(ctx, args...)
		}

		if err != nil {
			return values.None, e.maybeError(core.SourceError(e.src, err))
		}

		return out, nil
	}
}

func (e *FunctionCallExpression) maybeError(err error) error {
	if !e.ignoreErrors {
		return err
	}

	return nil
}
