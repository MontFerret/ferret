package operators

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type (
	UnaryOperatorType string
	UnaryOperator     struct {
		*baseOperator
		fn OperatorFunc
	}
)

const (
	UnaryOperatorTypeNoop     UnaryOperatorType = ""
	UnaryOperatorTypeNot      UnaryOperatorType = "!"
	UnaryOperatorTypeNot2     UnaryOperatorType = "NOT"
	UnaryOperatorTypeNegative UnaryOperatorType = "-"
	UnaryOperatorTypePositive UnaryOperatorType = "+"
)

var unaryOperators = map[UnaryOperatorType]OperatorFunc{
	UnaryOperatorTypeNoop:     ToBoolean,
	UnaryOperatorTypeNot:      Not,
	UnaryOperatorTypeNot2:     Not,
	UnaryOperatorTypeNegative: Negative,
	UnaryOperatorTypePositive: Positive,
}

func NewUnaryOperator(
	src core.SourceMap,
	exp core.Expression,
	operator UnaryOperatorType,
) (*UnaryOperator, error) {
	fn, exists := unaryOperators[operator]

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
