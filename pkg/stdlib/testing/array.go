package testing

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/MontFerret/ferret/pkg/stdlib/testing/base"
)

// ARRAY asserts that value is a array type.
// @param {Any} actual - Second to test.
// @param {String} [message] - Message to display on error.
var Array = base.Assertion{
	DefaultMessage: func(args []runtime.Value) string {
		return "be array"
	},
	MinArgs: 1,
	MaxArgs: 2,
	Fn: func(ctx context.Context, args []runtime.Value) (bool, error) {
		if err := runtime.AssertList(args[0]); err != nil {
			return false, err
		}

		return true, nil
	},
}
