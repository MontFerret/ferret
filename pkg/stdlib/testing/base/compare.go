package base

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type CompareOperator int

const (
	NotEqualOp       CompareOperator = 0
	EqualOp          CompareOperator = 1
	LessOp           CompareOperator = 2
	LessOrEqualOp    CompareOperator = 3
	GreaterOp        CompareOperator = 4
	GreaterOrEqualOp CompareOperator = 5
)

func (op CompareOperator) String() string {
	switch op {
	case NotEqualOp:
		return "not equal to"
	case EqualOp:
		return "equal to"
	case LessOp:
		return "less than"
	case LessOrEqualOp:
		return "less than or equal to"
	case GreaterOp:
		return "greater than"
	default:
		return "greater than or equal to"
	}
}

func (op CompareOperator) Compare(ctx context.Context, args []runtime.Value) (bool, error) {
	err := runtime.ValidateArgs(args, 2, 3)

	if err != nil {
		return false, err
	}

	actual := args[0]
	expected := args[1]

	switch op {
	case NotEqualOp:
		equal, err := runtime.EqualValues(ctx, actual, expected)
		return !bool(equal), err
	case EqualOp:
		equal, err := runtime.EqualValues(ctx, actual, expected)
		return bool(equal), err
	case LessOp:
		result, err := runtime.CompareValues(ctx, actual, expected)
		return result < 0, err
	case LessOrEqualOp:
		result, err := runtime.CompareValues(ctx, actual, expected)
		return result <= 0, err
	case GreaterOp:
		result, err := runtime.CompareValues(ctx, actual, expected)
		return result > 0, err
	default:
		result, err := runtime.CompareValues(ctx, actual, expected)
		return result >= 0, err
	}
}
