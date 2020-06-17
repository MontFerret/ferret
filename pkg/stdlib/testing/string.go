package testing

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
	"github.com/MontFerret/ferret/pkg/stdlib/testing/base"
)

// String asserts that value is a string type.
// @params actual (Mixed) - Value to test.
// @params message (String, optional) - Message to display on error.
var String = base.Assertion{
	DefaultMessage: func(args []core.Value) string {
		return "be string"
	},
	MinArgs: 1,
	MaxArgs: 2,
	Fn: func(ctx context.Context, args []core.Value) (bool, error) {
		return args[0].Type() == types.String, nil
	},
}
