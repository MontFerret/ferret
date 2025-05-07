package testing

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/stdlib/testing/base"
)

// EMPTY asserts that the target does not contain any values.
// @param {Measurable | Binary | Object | Any[] | String} actual - Second to test.
// @param {String} [message] - Message to display on error.
var Empty = base.Assertion{
	DefaultMessage: func(_ []runtime.Value) string {
		return "be empty"
	},
	MinArgs: 1,
	MaxArgs: 2,
	Fn: func(ctx context.Context, args []runtime.Value) (bool, error) {
		return runtime.IsEmpty(ctx, args[0])
	},
}
