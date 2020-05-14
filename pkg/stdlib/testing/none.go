package testing

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// None asserts that value is none.
// @param (Mixed) - Value to test.
// @param (String) - Message to display on error.
var None = Assertion{
	DefaultMessage: func(args []core.Value) string {
		return "be none"
	},
	MinArgs: 1,
	MaxArgs: 2,
	Fn: func(ctx context.Context, args []core.Value) (bool, error) {
		return args[0] == values.None, nil
	},
}
