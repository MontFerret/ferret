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
		err := runtime.AssertMap(args[0])
		return err == nil, nil
	},
}
