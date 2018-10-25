package operators

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type (
	EqualityOperator struct {
		*baseOperator
		fn OperatorFunc
	}
)

var equalityOperators = map[string]OperatorFunc{
	"==": Equal,
	"!=": NotEqual,
	">":  Greater,
	"<":  Less,
	">=": GreaterOrEqual,
	"<=": LessOrEqual,
}

func NewEqualityOperator(
	src core.SourceMap,
	left core.Expression,
	right core.Expression,
	operator string,
) (*EqualityOperator, error) {
	fn, exists := equalityOperators[operator]

	if !exists {
		return nil, core.Error(core.ErrInvalidArgument, "operator")
	}

	return &EqualityOperator{
		&baseOperator{src, left, right},
		fn,
	}, nil
}

func (operator *EqualityOperator) Exec(ctx context.Context, scope *core.Scope) (core.Value, error) {
	left, err := operator.left.Exec(ctx, scope)

	if err != nil {
		return nil, err
	}

	right, err := operator.right.Exec(ctx, scope)

	if err != nil {
		return nil, err
	}

	return operator.Eval(ctx, left, right)
}

func (operator *EqualityOperator) Eval(_ context.Context, left, right core.Value) (core.Value, error) {
	return operator.fn(left, right), nil
}
