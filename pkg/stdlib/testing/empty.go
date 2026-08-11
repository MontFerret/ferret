package testing

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/testing/base"
)

// Tests whether the target is empty.
// @param actual {Measurable | Binary | Object | Any[] | String} Value to test.
// @param message {String} Message to display on error.
// @return {Boolean} True when the configured assertion succeeds; otherwise an assertion error is returned.
var Empty = base.Assertion{
	DefaultMessage: func(_ []runtime.Value) string {
		return "be empty"
	},
	Args: base.Args{
		Min: 1,
		Max: 2,
	},
	Fn: func(ctx context.Context, args []runtime.Value) (bool, error) {
		value := args[0]

		// Validate that the value implements Measurable interface
		if err := runtime.AssertMeasurable(value); err != nil {
			return false, err
		}

		return runtime.IsEmpty(ctx, value)
	},
}
