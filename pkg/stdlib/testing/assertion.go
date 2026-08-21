package testing

import (
	"context"
	"errors"
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	assertionFn func(context.Context, []runtime.Value) (bool, error)

	messageFn func([]runtime.Value) string

	assertionArgs struct {
		min int
		max int
	}

	assertion struct {
		defaultMessage messageFn
		fn             assertionFn
		args           assertionArgs
	}
)

var errAssertion = errors.New("assertion error")

// fail returns an error.
// @param message {String} Message to display on error.
// @return {None} No success value is produced because this assertion always fails.
var failAssertion = assertion{
	defaultMessage: func(_ []runtime.Value) string {
		return "not fail"
	},
	args: assertionArgs{
		min: 0,
		max: 1,
	},
	fn: func(_ context.Context, _ []runtime.Value) (bool, error) {
		return false, nil
	},
}

func newTypeAssertion(expected runtime.Type, customExpectation ...string) assertion {
	if expected == nil {
		panic("unsupported type assertion for <nil>")
	}

	expectation := expected.String()
	if len(customExpectation) > 0 {
		expectation = customExpectation[0]
	}

	if expectation == "" {
		panic("type assertion expectation cannot be empty")
	}

	return assertion{
		defaultMessage: func(_ []runtime.Value) string {
			return fmt.Sprintf("be %s", expectation)
		},
		args: assertionArgs{
			min: 1,
			max: 2,
		},
		fn: func(_ context.Context, args []runtime.Value) (bool, error) {
			return expected.Is(args[0]), nil
		},
	}
}

func newComparisonAssertion(op comparisonOperator) assertion {
	return assertion{
		defaultMessage: func(args []runtime.Value) string {
			return fmt.Sprintf("be %s %s", op, formatValue(args[1]))
		},
		args: assertionArgs{
			min: 2,
			max: 3,
		},
		fn: func(ctx context.Context, args []runtime.Value) (bool, error) {
			return op.compare(ctx, args)
		},
	}
}

func (a assertion) positive() runtime.Function {
	return a.function(true)
}

func (a assertion) negative() runtime.Function {
	return a.function(false)
}

func (a assertion) function(positive bool) runtime.Function {
	return func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
		err := runtime.ValidateArgs(args, a.args.min, a.args.max)
		if err != nil {
			return runtime.None, err
		}

		result, err := a.fn(ctx, args)
		if err != nil {
			return runtime.None, err
		}

		if result == positive {
			return runtime.None, nil
		}

		return runtime.None, a.failure(args, positive)
	}
}

func (a assertion) failure(args []runtime.Value, positive bool) error {
	maxArgs := a.args.max
	if len(args) != maxArgs {
		connotation := ""
		if !positive {
			connotation = "not "
		}

		if maxArgs > 1 {
			actual := args[0]

			var message string
			if a.defaultMessage != nil {
				message = a.defaultMessage(args)
			} else if len(args) > 1 {
				message = fmt.Sprintf("be %s", args[1].String())
			} else {
				message = "exist"
			}

			return runtime.Error(errAssertion, fmt.Sprintf("expected %s %sto %s", formatValue(actual), connotation, message))
		}

		return runtime.Error(errAssertion, fmt.Sprintf("expected to %s%s", connotation, a.defaultMessage(args)))
	}

	// The last accepted argument is always a caller-provided message.
	message := args[maxArgs-1]

	return runtime.Error(errAssertion, message.String())
}
