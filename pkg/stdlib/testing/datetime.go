package testing

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/values"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/stdlib/testing/base"
)

// DATETIME asserts that value is a datetime type.
// @param {Any} actual - Value to test.
// @param {String} [message] - Message to display on error.
var DateTime = base.Assertion{
	DefaultMessage: func(args []core.Value) string {
		return "be datetime"
	},
	MinArgs: 1,
	MaxArgs: 2,
	Fn: func(ctx context.Context, args []core.Value) (bool, error) {
		if err := values.AssertDateTime(args[0]); err != nil {
			return false, err
		}

		return true, nil
	},
}
