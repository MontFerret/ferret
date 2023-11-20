package testing

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/values"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/stdlib/testing/base"
)

// STRING asserts that value is a string type.
// @param {Any} actual - Value to test.
// @param {String} [message] - Message to display on error.
var String = base.Assertion{
	DefaultMessage: func(args []core.Value) string {
		return "be string"
	},
	MinArgs: 1,
	MaxArgs: 2,
	Fn: func(ctx context.Context, args []core.Value) (bool, error) {
		if err := values.AssertString(args[0]); err != nil {
			return false, err
		}

		return true, nil
	},
}
