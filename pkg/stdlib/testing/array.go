package testing

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/stdlib/testing/base"
)

// ARRAY asserts that value is a array type.
// @param {Any} actual - Second to test.
// @param {String} [message] - Message to display on error.
var Array = base.Assertion{
	DefaultMessage: func(args []core.Value) string {
		return "be array"
	},
	MinArgs: 1,
	MaxArgs: 2,
	Fn: func(ctx context.Context, args []core.Value) (bool, error) {
		if err := core.AssertList(args[0]); err != nil {
			return false, err
		}

		return true, nil
	},
}
