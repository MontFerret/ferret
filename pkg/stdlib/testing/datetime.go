package testing

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
	"github.com/MontFerret/ferret/pkg/stdlib/testing/base"
)

// DateTime asserts that value is a datetime type.
// @param actual {Value} - Value to test.
// @param message {String, optional} - Message to display on error.
var DateTime = base.Assertion{
	DefaultMessage: func(args []core.Value) string {
		return "be datetime"
	},
	MinArgs: 1,
	MaxArgs: 2,
	Fn: func(ctx context.Context, args []core.Value) (bool, error) {
		return args[0].Type() == types.DateTime, nil
	},
}
