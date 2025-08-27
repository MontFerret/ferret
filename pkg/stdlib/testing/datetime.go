package testing

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/MontFerret/ferret/pkg/stdlib/testing/base"
)

// DATETIME asserts that value is a datetime type.
// @param {Any} actual - Second to test.
// @param {String} [message] - Message to display on error.
var DateTime = base.Assertion{
	DefaultMessage: func(args []runtime.Value) string {
		return "be datetime"
	},
	MinArgs: 1,
	MaxArgs: 2,
	Fn: func(ctx context.Context, args []runtime.Value) (bool, error) {
		// Check if the argument is a datetime type using CastDateTime
		// If casting succeeds, it's a datetime; if it fails, it's not
		_, err := runtime.CastDateTime(args[0])
		return err == nil, nil
	},
}
