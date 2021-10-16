package operators

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type (
	UnaryOperatorVariant string

	UnaryOperator struct {
		*baseOperator
		fn OperatorFunc
	}
)

const (
	UnaryOperatorVariantNoop     UnaryOperatorVariant = ""
	UnaryOperatorVariantNot      UnaryOperatorVariant = "!"
	UnaryOperatorVariantNot2     UnaryOperatorVariant = "NOT"
	UnaryOperatorVariantNegative UnaryOperatorVariant = "-"
	UnaryOperatorVariantPositive UnaryOperatorVariant = "+"
)

var unaryOperatorVariants = map[UnaryOperatorVariant]OperatorFunc{
	UnaryOperatorVariantNoop:     ToBoolean,
	UnaryOperatorVariantNot:      Not,
	UnaryOperatorVariantNot2:     Not,
	UnaryOperatorVariantNegative: Negative,
	UnaryOperatorVariantPositive: Positive,
}

func NewUnaryOperator(
	src core.SourceMap,
	exp core.Expression,
	variantStr string,
) (*UnaryOperator, error) {
	variant := UnaryOperatorVariant(variantStr)
	fn, exists := unaryOperatorVariants[variant]

	if !exists {
		return nil, core.Error(core.ErrInvalidArgument, "operator")
	}

	return &UnaryOperator{
		&baseOperator{
			src,
			exp,
			nil,
		},
		fn,
	}, nil
}

func (operator *UnaryOperator) Exec(ctx context.Context, scope *core.Scope) (core.Value, error) {
	value, err := operator.left.Exec(ctx, scope)

	if err != nil {
		return nil, core.SourceError(operator.src, err)
	}

	return operator.Eval(ctx, value, nil)
}

func (operator *UnaryOperator) Eval(_ context.Context, left, _ core.Value) (core.Value, error) {
	return operator.fn(left, values.None), nil
}
