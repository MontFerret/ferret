package testing

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// False asserts that value is false.
// @param (Mixed) - Value to test.
// @param (String) - Message to display on error.
var False = Assertion{
	DefaultMessage: func(args []core.Value) string {
		return "be false"
	},
	MinArgs: 1,
	MaxArgs: 2,
	Fn: func(ctx context.Context, args []core.Value) (bool, error) {
		return args[0] == values.False, nil
	},
}
