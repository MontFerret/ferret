package testing

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
	"github.com/MontFerret/ferret/pkg/stdlib/testing/base"
)

// FLOAT asserts that value is a float type.
// @param {Any} actual - Value to test.
// @param {String} [message] - Message to display on error.
var Float = base.Assertion{
	DefaultMessage: func(args []core.Value) string {
		return "be float"
	},
	MinArgs: 1,
	MaxArgs: 2,
	Fn: func(ctx context.Context, args []core.Value) (bool, error) {
		return args[0].Type() == types.Float, nil
	},
}
