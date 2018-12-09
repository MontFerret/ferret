package operators

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type (
	MathOperatorType string

	MathOperator struct {
		*baseOperator
		opType   MathOperatorType
		fn       OperatorFunc
		leftOnly bool
	}
)

const (
	MathOperatorTypeAdd       MathOperatorType = "+"
	MathOperatorTypeSubtract  MathOperatorType = "-"
	MathOperatorTypeMultiply  MathOperatorType = "*"
	MathOperatorTypeDivide    MathOperatorType = "/"
	MathOperatorTypeModulus   MathOperatorType = "%"
	MathOperatorTypeIncrement MathOperatorType = "++"
	MathOperatorTypeDecrement MathOperatorType = "--"
)

var mathOperators = map[MathOperatorType]OperatorFunc{
	MathOperatorTypeAdd:       Add,
	MathOperatorTypeSubtract:  Subtract,
	MathOperatorTypeMultiply:  Multiply,
	MathOperatorTypeDivide:    Divide,
	MathOperatorTypeModulus:   Modulus,
	MathOperatorTypeIncrement: Increment,
	MathOperatorTypeDecrement: Decrement,
}

func NewMathOperator(
	src core.SourceMap,
	left core.Expression,
	right core.Expression,
	operator MathOperatorType,
) (*MathOperator, error) {
	fn, exists := mathOperators[operator]

	if !exists {
		return nil, core.Error(core.ErrInvalidArgument, "operator type")
	}

	var leftOnly bool

	if operator == "++" || operator == "--" {
		leftOnly = true
	}

	return &MathOperator{
		&baseOperator{src, left, right},
		operator,
		fn,
		leftOnly,
	}, nil
}

func (operator *MathOperator) Type() MathOperatorType {
	return operator.opType
}

func (operator *MathOperator) Exec(ctx context.Context, scope *core.Scope) (core.Value, error) {
	left, err := operator.left.Exec(ctx, scope)

	if err != nil {
		return nil, err
	}

	if operator.leftOnly {
		return operator.Eval(ctx, left, values.None)
	}

	right, err := operator.right.Exec(ctx, scope)

	if err != nil {
		return nil, err
	}

	return operator.Eval(ctx, left, right)
}

func (operator *MathOperator) Eval(_ context.Context, left, right core.Value) (core.Value, error) {
	if operator.leftOnly {
		return operator.fn(left, values.None), nil
	}

	return operator.fn(left, right), nil
}
