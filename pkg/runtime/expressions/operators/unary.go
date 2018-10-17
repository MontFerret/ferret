package operators

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type (
	UnaryOperatorType int
	UnaryOperator     struct {
		*baseOperator
		value UnaryOperatorType
	}
)

const (
	UnaryOperatorTypeNot UnaryOperatorType = 0
	UnaryOperatorTypeYes UnaryOperatorType = 1
)

var unaryOperators = map[string]UnaryOperatorType{
	"!":   UnaryOperatorTypeNot,
	"NOT": UnaryOperatorTypeNot,
	"":    UnaryOperatorTypeYes,
}

func NewUnaryOperator(
	src core.SourceMap,
	exp core.Expression,
	operator string,
) (*UnaryOperator, error) {
	op, exists := unaryOperators[operator]

	if !exists {
		return nil, core.Error(core.ErrInvalidArgument, "operator")
	}

	return &UnaryOperator{
		&baseOperator{
			src,
			exp,
			nil,
		},
		op,
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
	if operator.value == UnaryOperatorTypeNot {
		return Not(left, values.None), nil
	}

	return values.ToBoolean(left), nil
}
