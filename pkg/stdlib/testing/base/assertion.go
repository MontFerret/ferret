package base

import (
	"context"
	"fmt"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type AssertionFn func(ctx context.Context, args []core.Value) (bool, error)

type MessageFn func(args []core.Value) string

type Assertion struct {
	DefaultMessage MessageFn
	MinArgs        int
	MaxArgs        int
	Fn             AssertionFn
}

func NewPositiveAssertion(assertion Assertion) core.Function {
	return newInternal(assertion, true)
}

func NewNegativeAssertion(assertion Assertion) core.Function {
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

			return core.Error(ErrAssertion, fmt.Sprintf("expected %s %sto %s", FormatValue(actual), connotation, msg))
		}

		return core.Error(ErrAssertion, fmt.Sprintf("expected to %s%s", connotation, assertion.DefaultMessage(args)))
	}

	// Last argument is always is a custom message
	msg := args[assertion.MaxArgs-1]

	return core.Error(ErrAssertion, msg.String())
}
