package testing

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/MontFerret/ferret/pkg/stdlib/testing/base"
)

// INT asserts that value is a int type.
// @param {Any} actual - Actual value.
// @param {String} [message] - Message to display on error.
var Int = base.Assertion{
	DefaultMessage: func(args []runtime.Value) string {
		return "be int"
	},
	MinArgs: 1,
	MaxArgs: 2,
	Fn: func(ctx context.Context, args []runtime.Value) (bool, error) {
		// Check if the argument is an int type using CastInt  
		// If casting succeeds, it's an int; if it fails, it's not
		_, err := runtime.CastInt(args[0])
		return err == nil, nil
	},
}
