package testing

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
	"github.com/MontFerret/ferret/pkg/stdlib/testing/base"
)

// INT asserts that value is a int type.
// @param {Any} actual - Actual value.
// @param {String} [message] - Message to display on error.
var Int = base.Assertion{
	DefaultMessage: func(args []core.Value) string {
		return "be int"
	},
	MinArgs: 1,
	MaxArgs: 2,
	Fn: func(ctx context.Context, args []core.Value) (bool, error) {
		return args[0].Type() == types.Int, nil
	},
}
