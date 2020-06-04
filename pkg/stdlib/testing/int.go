package testing

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// Int asserts that value is a int type.
// @param (Mixed) - Value to test.
// @param (String) - Message to display on error.
var Int = Assertion{
	DefaultMessage: func(args []core.Value) string {
		return "be int"
	},
	MinArgs: 1,
	MaxArgs: 2,
	Fn: func(ctx context.Context, args []core.Value) (bool, error) {
		return args[0].Type() == types.Int, nil
	},
}
