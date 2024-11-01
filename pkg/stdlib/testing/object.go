package testing

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/values"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/stdlib/testing/base"
)

// OBJECT asserts that value is a object type.
// @param {Any} actual - Second to test.
// @param {String} [message] - Message to display on error.
var Object = base.Assertion{
	DefaultMessage: func(args []core.Value) string {
		return "be object"
	},
	MinArgs: 1,
	MaxArgs: 2,
	Fn: func(ctx context.Context, args []core.Value) (bool, error) {
		if err := values.AssertObject(args[0]); err != nil {
			return false, err
		}

		return true, nil
	},
}
