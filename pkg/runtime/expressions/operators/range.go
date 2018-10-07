package operators

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type RangeOperator struct {
	*baseOperator
}

func NewRangeOperator(
	src core.SourceMap,
	left core.Expression,
	right core.Expression,
) (*RangeOperator, error) {
	if core.IsNil(left) {
		return nil, core.Error(core.ErrMissedArgument, "left expression")
	}

	if core.IsNil(right) {
		return nil, core.Error(core.ErrMissedArgument, "right expression")
	}

	return &RangeOperator{&baseOperator{src, left, right}}, nil
}

func (operator *RangeOperator) Iterate(ctx context.Context, scope *core.Scope) (collections.Iterator, error) {
	arr, err := operator.Exec(ctx, scope)

	if err != nil {
		return nil, err
	}

	return collections.NewArrayIterator(arr.(*values.Array)), nil
}

func (operator *RangeOperator) Exec(ctx context.Context, scope *core.Scope) (core.Value, error) {
	left, err := operator.left.Exec(ctx, scope)

	if err != nil {
		return values.None, core.SourceError(operator.src, err)
	}

	right, err := operator.right.Exec(ctx, scope)

	if err != nil {
		return values.None, core.SourceError(operator.src, err)
	}

	return operator.Eval(ctx, left, right)
}

func (operator *RangeOperator) Eval(_ context.Context, left, right core.Value) (core.Value, error) {
	err := core.ValidateType(left, core.IntType)

	if err != nil {
		return values.None, core.SourceError(operator.src, err)
	}

	err = core.ValidateType(right, core.IntType)

	if err != nil {
		return values.None, core.SourceError(operator.src, err)
	}

	arr := values.NewArray(10)
	start := left.(values.Int)
	end := right.(values.Int)

	for i := start; i <= end; i++ {
		arr.Push(i)
	}

	return arr, nil
}
