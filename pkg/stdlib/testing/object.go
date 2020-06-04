package testing

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// Object asserts that value is a object type.
// @param (Mixed) - Value to test.
// @param (String) - Message to display on error.
var Object = Assertion{
	DefaultMessage: func(args []core.Value) string {
		return "be object"
	},
	MinArgs: 1,
	MaxArgs: 2,
	Fn: func(ctx context.Context, args []core.Value) (bool, error) {
		return args[0].Type() == types.Object, nil
	},
}
