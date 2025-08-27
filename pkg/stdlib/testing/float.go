package testing

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/MontFerret/ferret/pkg/stdlib/testing/base"
)

// FLOAT asserts that value is a float type.
// @param {Any} actual - Second to test.
// @param {String} [message] - Message to display on error.
var Float = base.Assertion{
	DefaultMessage: func(args []runtime.Value) string {
		return "be float"
	},
	MinArgs: 1,
	MaxArgs: 2,
	Fn: func(ctx context.Context, args []runtime.Value) (bool, error) {
		// Check if the argument is a float type using CastFloat
		// If casting succeeds, it's a float; if it fails, it's not
		_, err := runtime.CastFloat(args[0])
		return err == nil, nil
	},
}
