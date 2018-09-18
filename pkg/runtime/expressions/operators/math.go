package operators

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type MathOperator struct {
	*baseOperator
	fn       Operator
	leftOnly bool
}

var mathOperators = map[string]Operator{
	"+":  Add,
	"-":  Subtract,
	"*":  Multiply,
	"/":  Divide,
	"%":  Modulus,
	"++": Increment,
	"--": Decrement,
}

func NewMathOperator(
	src core.SourceMap,
	left core.Expression,
	right core.Expression,
	operator string,
) (*MathOperator, error) {
	fn, exists := mathOperators[operator]

	if !exists {
		return nil, core.Error(core.ErrInvalidArgument, "operator")
	}

	var leftOnly bool

	if operator == "++" || operator == "--" {
		leftOnly = true
	}

	return &MathOperator{
		&baseOperator{src, left, right},
		fn,
		leftOnly,
	}, nil
}

func (operator *MathOperator) Exec(ctx context.Context, scope *core.Scope) (core.Value, error) {
	left, err := operator.left.Exec(ctx, scope)

	if err != nil {
		return nil, err
	}

	if operator.leftOnly {
		return operator.fn(left, values.None), nil
	}

	right, err := operator.right.Exec(ctx, scope)

	if err != nil {
		return nil, err
	}

	return operator.fn(left, right), nil
}
