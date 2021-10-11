package operators

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type (
	MathOperatorVariant string

	MathOperator struct {
		*baseOperator
		variant  MathOperatorVariant
		fn       OperatorFunc
		leftOnly bool
	}
)

const (
	MathOperatorVariantAdd       MathOperatorVariant = "+"
	MathOperatorVariantSubtract  MathOperatorVariant = "-"
	MathOperatorVariantMultiply  MathOperatorVariant = "*"
	MathOperatorVariantDivide    MathOperatorVariant = "/"
	MathOperatorVariantModulus   MathOperatorVariant = "%"
	MathOperatorVariantIncrement MathOperatorVariant = "++"
	MathOperatorVariantDecrement MathOperatorVariant = "--"
)

var mathOperatorVariants = map[MathOperatorVariant]OperatorFunc{
	MathOperatorVariantAdd:       Add,
	MathOperatorVariantSubtract:  Subtract,
	MathOperatorVariantMultiply:  Multiply,
	MathOperatorVariantDivide:    Divide,
	MathOperatorVariantModulus:   Modulus,
	MathOperatorVariantIncrement: Increment,
	MathOperatorVariantDecrement: Decrement,
}

func NewMathOperator(
	src core.SourceMap,
	left core.Expression,
	right core.Expression,
	variantStr string,
) (*MathOperator, error) {
	variant := MathOperatorVariant(variantStr)
	fn, exists := mathOperatorVariants[variant]

	if !exists {
		return nil, core.Error(core.ErrInvalidArgument, "operator type")
	}

	var leftOnly bool

	if variant == "++" || variant == "--" {
		leftOnly = true
	}

	return &MathOperator{
		&baseOperator{src, left, right},
		variant,
		fn,
		leftOnly,
	}, nil
}

func (operator *MathOperator) Type() MathOperatorVariant {
	return operator.variant
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
