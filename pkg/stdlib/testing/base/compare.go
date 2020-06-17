package base

import "github.com/MontFerret/ferret/pkg/runtime/core"

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
		result = actual.Compare(expected) != 0
	case EqualOp:
		result = actual.Compare(expected) == 0
	case LessOp:
		result = actual.Compare(expected) == -1
	case LessOrEqualOp:
		c := actual.Compare(expected)
		result = c == -1 || c == 0
	case GreaterOp:
		result = actual.Compare(expected) == 1
	default:
		c := actual.Compare(expected)
		result = c == 1 || c == 0
	}

	return result, nil
}
