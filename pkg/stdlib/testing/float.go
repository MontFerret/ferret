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
		err := runtime.AssertFloat(args[0])
		return err == nil, nil
	},
}
