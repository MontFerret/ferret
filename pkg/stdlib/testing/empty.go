package testing

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/stdlib/testing/base"
)

// EMPTY asserts that the target does not contain any values.
// @param {Measurable | Binary | Object | Any[] | String} actual - Second to test.
// @param {String} [message] - Message to display on error.
var Empty = base.Assertion{
	DefaultMessage: func(_ []runtime.Value) string {
		return "be empty"
	},
	MinArgs: 1,
	MaxArgs: 2,
	Fn: func(ctx context.Context, args []runtime.Value) (bool, error) {
		value := args[0]
		
		// Validate that the value is one of the supported types
		if err := runtime.AssertString(value); err != nil {
			if err := runtime.AssertList(value); err != nil {
				if err := runtime.AssertMap(value); err != nil {
					if err := runtime.AssertBinary(value); err != nil {
						// If none of the supported types match, return the last error
						return false, err
					}
				}
			}
		}
		
		return runtime.IsEmpty(ctx, value)
	},
}
