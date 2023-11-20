package base

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
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

func (op CompareOperator) Compare(args []core.Value) (bool, error) {
	err := core.ValidateArgs(args, 2, 3)

	if err != nil {
		return false, err
	}

	actual := args[0]
	expected := args[1]

	var result bool

	switch op {
	case NotEqualOp:
		result = values.Compare(actual, expected) != 0
	case EqualOp:
		result = values.Compare(actual, expected) == 0
	case LessOp:
		result = values.Compare(actual, expected) == -1
	case LessOrEqualOp:
		c := values.Compare(actual, expected)
		result = c == -1 || c == 0
	case GreaterOp:
		result = values.Compare(actual, expected) == 1
	default:
		c := values.Compare(actual, expected)
		result = c == 1 || c == 0
	}

	return result, nil
}
