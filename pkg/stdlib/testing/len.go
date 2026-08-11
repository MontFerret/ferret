package testing

import (
	"context"
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/testing/base"
)

// Tests whether a measurable value has the expected length or size.
// @param actual {Measurable} Measurable value.
// @param length {Int} Target length.
// @param message {String} Message to display on error.
// @return {Boolean} True when the configured assertion succeeds; otherwise an assertion error is returned.
var Len = base.Assertion{
	DefaultMessage: func(args []runtime.Value) string {
		return fmt.Sprintf("has size %s", args[1])
	},
	Args: base.Args{
		Min: 2,
		Max: 3,
	},
	Fn: func(ctx context.Context, args []runtime.Value) (bool, error) {
		col := args[0]
		size := args[1]

		// Validate that the value implements Measurable interface
		if err := runtime.AssertMeasurable(col); err != nil {
			return false, err
		}

		out, err := runtime.Length(ctx, col)

		if err != nil {
			return false, err
		}

		equal, err := runtime.EqualValues(ctx, out, size)
		return bool(equal), err
	},
}
