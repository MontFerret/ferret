package operators

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

type InOperator struct {
	*baseOperator
	negate bool
}

func NewInOperator(
	src core.SourceMap,
	left core.Expression,
	right core.Expression,
	negate bool,
) (*InOperator, error) {
	if left == nil {
		return nil, core.Error(core.ErrMissedArgument, "left expression")
	}

	if right == nil {
		return nil, core.Error(core.ErrMissedArgument, "right expression")
	}

	return &InOperator{&baseOperator{src, left, right}, negate}, nil
}

func (operator *InOperator) Exec(ctx context.Context, scope *core.Scope) (core.Value, error) {
	left, err := operator.left.Exec(ctx, scope)

	if err != nil {
		return values.False, core.SourceError(operator.src, err)
	}

	right, err := operator.right.Exec(ctx, scope)

	if err != nil {
		return values.False, core.SourceError(operator.src, err)
	}

	return operator.Eval(ctx, left, right)
}

func (operator *InOperator) Eval(_ context.Context, left, right core.Value) (core.Value, error) {
	err := core.ValidateType(right, types.Array)

	if err != nil {
		// TODO: Return the error? AQL just returns false
		return values.False, nil
	}

	arr := right.(*values.Array)
	found := arr.IndexOf(left) > -1

	if operator.negate {
		return values.NewBoolean(!found), nil
	}

	return values.NewBoolean(found), nil
}
