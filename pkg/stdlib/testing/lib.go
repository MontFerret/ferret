package testing

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

var (
	ErrAssertion = errors.New("assertion error")
)

type AssertionFn func(ctx context.Context, args []core.Value) (bool, error)

type MessageFn func(args []core.Value) string

type Assertion struct {
	DefaultMessage MessageFn
	MinArgs        int
	MaxArgs        int
	Fn             AssertionFn
}

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

func RegisterLib(ns core.Namespace) error {
	t := ns.Namespace("T")

	if err := registerNOT(t); err != nil {
		return err
	}

	return t.RegisterFunctions(
		core.NewFunctionsFromMap(map[string]core.Function{
			"EMPTY":   NewPositive(Empty),
			"EQUAL":   NewPositive(Equal),
			"FAIL":    NewPositive(Fail),
			"FALSE":   NewPositive(False),
			"GT":      NewPositive(Gt),
			"GTE":     NewPositive(Gte),
			"INCLUDE": NewPositive(Include),
			"LEN":     NewPositive(Len),
			"MATCH":   NewPositive(Match),
			"LT":      NewPositive(Lt),
			"LTE":     NewPositive(Lte),
			"NONE":    NewPositive(None),
			"TRUE":    NewPositive(True),
		}),
	)
}

func registerNOT(ns core.Namespace) error {
	t := ns.Namespace("NOT")

	return t.RegisterFunctions(
		core.NewFunctionsFromMap(map[string]core.Function{
			"EMPTY":   NewNegative(Empty),
			"EQUAL":   NewNegative(Equal),
			"FALSE":   NewNegative(False),
			"GT":      NewNegative(Gt),
			"GTE":     NewNegative(Gte),
			"INCLUDE": NewNegative(Include),
			"LEN":     NewNegative(Len),
			"MATCH":   NewNegative(Match),
			"LT":      NewNegative(Lt),
			"LTE":     NewNegative(Lte),
			"NONE":    NewNegative(None),
			"TRUE":    NewNegative(True),
		}),
	)
}

func compare(args []core.Value, op CompareOperator) (bool, error) {
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

func NewPositive(assertion Assertion) core.Function {
	return newInternal(assertion, true)
}

func NewNegative(assertion Assertion) core.Function {
	return newInternal(assertion, false)
}

func newInternal(assertion Assertion, connotation bool) core.Function {
	return func(ctx context.Context, args ...core.Value) (core.Value, error) {
		err := core.ValidateArgs(args, assertion.MinArgs, assertion.MaxArgs)

		if err != nil {
			return values.None, err
		}

		res, err := assertion.Fn(ctx, args)

		if err != nil {
			return values.None, err
		}

		if res == connotation {
			return values.None, nil
		}

		return values.None, toError(assertion, args, connotation)
	}
}

func toError(assertion Assertion, args []core.Value, positive bool) error {
	if len(args) != assertion.MaxArgs {
		connotation := ""

		if !positive {
			connotation = "not "
		}

		if assertion.MaxArgs > 1 {
			actual := args[0]

			return core.Error(ErrAssertion, fmt.Sprintf("expected %s %sto %s", formatValue(actual), connotation, assertion.DefaultMessage(args)))
		}

		return core.Error(ErrAssertion, fmt.Sprintf("expected to %s%s", connotation, assertion.DefaultMessage(args)))
	}

	// Last argument is always is a custom message
	msg := args[assertion.MaxArgs-1]

	return core.Error(ErrAssertion, msg.String())
}

func formatValue(val core.Value) string {
	valStr := val.String()

	if val == values.None {
		valStr = "none"
	}

	return fmt.Sprintf("[%s] '%s'", val.Type(), valStr)
}
