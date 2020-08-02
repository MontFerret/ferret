package testing

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
	"github.com/MontFerret/ferret/pkg/stdlib/testing/base"
)

// Binary asserts that value is a binary type.
// @param actual {Any} - Value to test.
// @param message {String, optional} - Message to display on error.
var Binary = base.Assertion{
	DefaultMessage: func(args []core.Value) string {
		return "be binary"
	},
	MinArgs: 1,
	MaxArgs: 2,
	Fn: func(ctx context.Context, args []core.Value) (bool, error) {
		return args[0].Type() == types.Binary, nil
	},
}
