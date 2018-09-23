package expressions

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type FunctionCallExpression struct {
	src  core.SourceMap
	fun  core.Function
	args []core.Expression
}

func NewFunctionCallExpression(
	src core.SourceMap,
	fun core.Function,
	args ...core.Expression,
) (*FunctionCallExpression, error) {
	if fun == nil {
		return nil, core.Error(core.ErrMissedArgument, "function")
	}

	return &FunctionCallExpression{src, fun, args}, nil
}

func (e *FunctionCallExpression) Iterate(ctx context.Context, scope *core.Scope) (collections.Iterator, error) {
	value, err := e.Exec(ctx, scope)

	if err != nil {
		return nil, core.SourceError(e.src, err)
	}

	iter, err := collections.ToIterator(value)

	if err != nil {
		return nil, core.SourceError(e.src, err)
	}

	return iter, nil
}

func (e *FunctionCallExpression) Exec(ctx context.Context, scope *core.Scope) (core.Value, error) {
	var out core.Value
	var err error

	if len(e.args) == 0 {
		out, err = e.fun(ctx)
	} else {
		args := make([]core.Value, len(e.args))

		for idx, arg := range e.args {
			out, err := arg.Exec(ctx, scope)

			if err != nil {
				return values.None, core.SourceError(e.src, err)
			}

			args[idx] = out
		}

		out, err = e.fun(ctx, args...)
	}

	if err != nil {
		return values.None, core.SourceError(e.src, err)
	}

	return out, nil
}
