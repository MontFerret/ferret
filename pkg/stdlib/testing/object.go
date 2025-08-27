package testing

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/MontFerret/ferret/pkg/stdlib/testing/base"
)

// OBJECT asserts that value is a object type.
// @param {Any} actual - Second to test.
// @param {String} [message] - Message to display on error.
var Object = base.Assertion{
	DefaultMessage: func(args []runtime.Value) string {
		return "be object"
	},
	MinArgs: 1,
	MaxArgs: 2,
	Fn: func(ctx context.Context, args []runtime.Value) (bool, error) {
		// Check if the argument is an object/map type using CastMap
		// If casting succeeds, it's a map; if it fails, it's not
		_, err := runtime.CastMap(args[0])
		return err == nil, nil
	},
}
