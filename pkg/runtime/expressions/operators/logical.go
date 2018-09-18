package operators

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type (
	LogicalOperatorType int
	LogicalOperator     struct {
		*baseOperator
		value LogicalOperatorType
	}
)

var (
	AndType LogicalOperatorType = 0
	OrType  LogicalOperatorType = 1
	NotType LogicalOperatorType = 2
)

var logicalOperators = map[string]LogicalOperatorType{
	"&&":  AndType,
	"AND": AndType,
	"||":  OrType,
	"OR":  OrType,
	"NOT": NotType,
}

func NewLogicalOperator(
	src core.SourceMap,
	left core.Expression,
	right core.Expression,
	operator string,
) (*LogicalOperator, error) {
	op, exists := logicalOperators[operator]

	if !exists {
		return nil, core.Error(core.ErrInvalidArgument, "value")
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
	left, err := operator.left.Exec(ctx, scope)

	if err != nil {
		return nil, err
	}

	left = operator.ensureType(left)

	if operator.value == NotType {
		return Not(left, values.None), nil
	}

	if operator.value == AndType && left == values.False {
		return values.False, nil
	}

	if operator.value == OrType && left == values.True {
		return values.True, nil
	}

	right, err := operator.right.Exec(ctx, scope)

	if err != nil {
		return nil, err
	}

	return operator.ensureType(right), nil
}

func (operator *LogicalOperator) ensureType(value core.Value) core.Value {
	if value.Type() != core.BooleanType {
		if value.Type() == core.NoneType {
			return values.False
		}

		return values.True
	}

	return value
}
