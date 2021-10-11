package operators

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

type (
	LogicalOperatorVariant int

	LogicalOperator struct {
		*baseOperator
		variant LogicalOperatorVariant
	}
)

const (
	LogicalOperatorVariantAnd LogicalOperatorVariant = 0
	LogicalOperatorVariantOr  LogicalOperatorVariant = 1
	LogicalOperatorVariantNot LogicalOperatorVariant = 2
)

var logicalOperators = map[string]LogicalOperatorVariant{
	"&&":  LogicalOperatorVariantAnd,
	"AND": LogicalOperatorVariantAnd,
	"||":  LogicalOperatorVariantOr,
	"OR":  LogicalOperatorVariantOr,
	"NOT": LogicalOperatorVariantNot,
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
	if operator.variant == LogicalOperatorVariantNot {
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

	if operator.variant == LogicalOperatorVariantAnd && leftBool == values.False {
		if left.Type() == types.Boolean {
			return values.False, nil
		}

		return left, nil
	}

	if operator.variant == LogicalOperatorVariantOr && leftBool == values.True {
		return left, nil
	}

	right, err := operator.right.Exec(ctx, scope)

	if err != nil {
		return values.None, core.SourceError(operator.src, err)
	}

	return right, nil
}

func (operator *LogicalOperator) Eval(_ context.Context, left, right core.Value) (core.Value, error) {
	if operator.variant == LogicalOperatorVariantNot {
		return Not(right, values.None), nil
	}

	leftBool := values.ToBoolean(left)

	if operator.variant == LogicalOperatorVariantAnd && leftBool == values.False {
		if left.Type() == types.Boolean {
			return values.False, nil
		}

		return left, nil
	}

	if operator.variant == LogicalOperatorVariantOr && leftBool == values.True {
		return left, nil
	}

	return right, nil
}
