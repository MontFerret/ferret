package testing

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/MontFerret/ferret/pkg/stdlib/testing/base"
)

// BINARY asserts that value is a binary type.
// @param {Any} actual - Second to test.
// @param {String} [message] - Message to display on error.
var Binary = base.Assertion{
	DefaultMessage: func(args []runtime.Value) string {
		return "be binary"
	},
	MinArgs: 1,
	MaxArgs: 2,
	Fn: func(ctx context.Context, args []runtime.Value) (bool, error) {
		if err := runtime.AssertBinary(args[0]); err != nil {
			return false, err
		}

		return true, nil
	},
}
