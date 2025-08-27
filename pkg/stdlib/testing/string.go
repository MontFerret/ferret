package testing

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/MontFerret/ferret/pkg/stdlib/testing/base"
)

// STRING asserts that value is a string type.
// @param {Any} actual - Second to test.
// @param {String} [message] - Message to display on error.
var String = base.Assertion{
	DefaultMessage: func(args []runtime.Value) string {
		return "be string"
	},
	MinArgs: 1,
	MaxArgs: 2,
	Fn: func(ctx context.Context, args []runtime.Value) (bool, error) {
		err := runtime.AssertString(args[0])
		return err == nil, nil
	},
}
