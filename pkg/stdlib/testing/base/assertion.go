package base

import (
	"context"
	"fmt"

	"github.com/MontFerret/ferret/pkg/runtime"
)

type AssertionFn func(ctx context.Context, args []runtime.Value) (bool, error)

type MessageFn func(args []runtime.Value) string

type Assertion struct {
	DefaultMessage MessageFn
	MinArgs        int
	MaxArgs        int
	Fn             AssertionFn
}

func NewPositiveAssertion(assertion Assertion) runtime.Function {
	return newInternal(assertion, true)
}

func NewNegativeAssertion(assertion Assertion) runtime.Function {
	return newInternal(assertion, false)
}

func newInternal(assertion Assertion, connotation bool) runtime.Function {
	return func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
		err := runtime.ValidateArgs(args, assertion.MinArgs, assertion.MaxArgs)

		if err != nil {
			return runtime.None, err
		}

		res, err := assertion.Fn(ctx, args)

		if err != nil {
			return runtime.None, err
		}

		if res == connotation {
			return runtime.None, nil
		}

		return runtime.None, toError(assertion, args, connotation)
	}
}

func toError(assertion Assertion, args []runtime.Value, positive bool) error {
	if len(args) != assertion.MaxArgs {
		connotation := ""

		if !positive {
			connotation = "not "
		}

		if assertion.MaxArgs > 1 {
			actual := args[0]

			var msg string

			if assertion.DefaultMessage != nil {
				msg = assertion.DefaultMessage(args)
			} else {
				if len(args) > 1 {
					msg = fmt.Sprintf("be %s", args[1].String())
				} else {
					msg = "exist"
				}
			}

			return runtime.Error(ErrAssertion, fmt.Sprintf("expected %s %sto %s", FormatValue(actual), connotation, msg))
		}

		return runtime.Error(ErrAssertion, fmt.Sprintf("expected to %s%s", connotation, assertion.DefaultMessage(args)))
	}

	// Last argument is always is a custom message
	msg := args[assertion.MaxArgs-1]

	return runtime.Error(ErrAssertion, msg.String())
}
