package testing

import (
	"context"

	"github.com/pkg/errors"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

var (
	ErrAssertion = errors.New("assertion error")
)

type AssertionFn func(ctx context.Context, args []core.Value) (values.Boolean, error)

type MessageFn func(args []core.Value) values.String

type Assertion struct {
	DefaultMessage MessageFn
	MessageArg     int
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
			"ASSERT":  positive(Assert),
			"EMPTY":   positive(Empty),
			"EQUAL":   positive(Equal),
			"FAIL":    positive(Fail),
			"FALSE":   positive(False),
			"GT":      positive(Gt),
			"GTE":     positive(Gte),
			"INCLUDE": positive(Include),
			"LEN":     positive(Len),
			"MATCH":   positive(Match),
			"LT":      positive(Lt),
			"LTE":     positive(Lte),
			"NONE":    positive(None),
			"TRUE":    positive(True),
		}),
	)
}

func registerNOT(ns core.Namespace) error {
	t := ns.Namespace("NOT")

	return t.RegisterFunctions(
		core.NewFunctionsFromMap(map[string]core.Function{
			"EMPTY":   negative(Empty),
			"EQUAL":   negative(Equal),
			"FAIL":    negative(Fail),
			"FALSE":   negative(False),
			"GT":      negative(Gt),
			"GTE":     negative(Gte),
			"INCLUDE": negative(Include),
			"LEN":     negative(Len),
			"MATCH":   negative(Match),
			"LT":      negative(Lt),
			"LTE":     negative(Lte),
			"NONE":    negative(None),
			"TRUE":    negative(True),
		}),
	)
}

func compare(args []core.Value, op CompareOperator) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 3)

	if err != nil {
		return values.None, err
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

	if result {
		return values.None, nil
	}

	if len(args) > 2 {
		return values.None, core.Error(ErrAssertion, args[2].String())
	}

	return values.None, core.Errorf(ErrAssertion, "expected %s to be %s %s", actual, op, expected)
}

func positive(assertion Assertion) core.Function {
	return func(ctx context.Context, args ...core.Value) (core.Value, error) {
		err := core.ValidateArgs(args, 1, 2)

		if err != nil {
			return values.None, err
		}

		res, err := assertion.Fn(ctx, args)

		if err != nil {
			return values.None, err
		}

		if res {
			return values.None, nil
		}

		if len(args) <= assertion.MessageArg {
			return values.None, core.Errorf(ErrAssertion, assertion.DefaultMessage(args).String())
		}

		msg := args[assertion.MessageArg]

		return values.None, core.Error(ErrAssertion, msg.String())
	}
}

func negative(assertion Assertion) core.Function {
	return func(ctx context.Context, args ...core.Value) (core.Value, error) {
		err := core.ValidateArgs(args, 1, 2)

		if err != nil {
			return values.None, err
		}

		res, err := assertion.Fn(ctx, args)

		if err != nil {
			return values.None, err
		}

		if !res {
			return values.None, nil
		}

		if len(args) <= assertion.MessageArg {
			return values.None, core.Error(ErrAssertion, assertion.DefaultMessage(args).String())
		}

		msg := args[assertion.MessageArg]

		return values.None, core.Error(ErrAssertion, msg.String())
	}
}
