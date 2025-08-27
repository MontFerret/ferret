package testing

import (
	"context"
	"fmt"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/stdlib/testing/base"
)

// LEN asserts that a measurable value has a length or size with the expected value.
// @param {Measurable} actual - Measurable value.
// @param {Int} length - Target length.
// @param {String} [message] - Message to display on error.
var Len = base.Assertion{
	DefaultMessage: func(args []runtime.Value) string {
		return fmt.Sprintf("has size %s", args[1])
	},
	MinArgs: 2,
	MaxArgs: 3,
	Fn: func(ctx context.Context, args []runtime.Value) (bool, error) {
		col := args[0]
		size := args[1]

		// Validate that the value is one of the supported measurable types
		if err := runtime.AssertString(col); err != nil {
			if err := runtime.AssertList(col); err != nil {
				if err := runtime.AssertMap(col); err != nil {
					if err := runtime.AssertBinary(col); err != nil {
						// If none of the supported types match, return the last error
						return false, err
					}
				}
			}
		}

		out, err := runtime.Length(ctx, col)

		if err != nil {
			return false, err
		}

		return runtime.CompareValues(out, size) == 0, nil
	},
}
