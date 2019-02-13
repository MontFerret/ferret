package operators

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

type (
	LogicalOperatorType int
	LogicalOperator     struct {
		*baseOperator
		value LogicalOperatorType
	}
)

const (
	LogicalOperatorTypeAnd LogicalOperatorType = 0
	LogicalOperatorTypeOr  LogicalOperatorType = 1
	LogicalOperatorTypeNot LogicalOperatorType = 2
)

var logicalOperators = map[string]LogicalOperatorType{
	"&&":  LogicalOperatorTypeAnd,
	"AND": LogicalOperatorTypeAnd,
	"||":  LogicalOperatorTypeOr,
	"OR":  LogicalOperatorTypeOr,
	"NOT": LogicalOperatorTypeNot,
}

func NewLogicalOperator(
	src core.SourceMap,
	left core.Expression,
	right core.Expression,
	operator string,
) (*LogicalOperator, error) {
	op, exists := logicalOperators[operator]

	if !exists {
		return nil, core.Error(core.ErrInvalidArgument, "operator")
	}

	return &LogicalOperator{
		&baseOperator{
			src,
			left,
			right,
		},
		op,
	}, nil
}

func (operator *LogicalOperator) Exec(ctx context.Context, scope *core.Scope) (core.Value, error) {
	if operator.value == LogicalOperatorTypeNot {
		val, err := operator.right.Exec(ctx, scope)

		if err != nil {
			return values.None, core.SourceError(operator.src, err)
		}

		return Not(val, values.None), nil
	}

	left, err := operator.left.Exec(ctx, scope)

	if err != nil {
		return values.None, core.SourceError(operator.src, err)
	}

	leftBool := values.ToBoolean(left)

	if operator.value == LogicalOperatorTypeAnd && leftBool == values.False {
		if left.Type() == types.Boolean {
			return values.False, nil
		}

		return left, nil
	}

	if operator.value == LogicalOperatorTypeOr && leftBool == values.True {
		return left, nil
	}

	right, err := operator.right.Exec(ctx, scope)

	if err != nil {
		return values.None, core.SourceError(operator.src, err)
	}

	return right, nil
}

func (operator *LogicalOperator) Eval(_ context.Context, left, right core.Value) (core.Value, error) {
	if operator.value == LogicalOperatorTypeNot {
		return Not(right, values.None), nil
	}

	leftBool := values.ToBoolean(left)

	if operator.value == LogicalOperatorTypeAnd && leftBool == values.False {
		if left.Type() == types.Boolean {
			return values.False, nil
		}

		return left, nil
	}

	if operator.value == LogicalOperatorTypeOr && leftBool == values.True {
		return left, nil
	}

	return right, nil
}
