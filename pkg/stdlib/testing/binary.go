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
		// Check if the argument is a binary type using CastBinary
		// If casting succeeds, it's a binary; if it fails, it's not
		_, err := runtime.CastBinary(args[0])
		return err == nil, nil
	},
}
