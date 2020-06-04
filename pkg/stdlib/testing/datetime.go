package testing

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// DateTime asserts that value is a datetime type.
// @params actual (Mixed) - Value to test.
// @params message (String, optional) - Message to display on error.
var DateTime = Assertion{
	DefaultMessage: func(args []core.Value) string {
		return "be datetime"
	},
	MinArgs: 1,
	MaxArgs: 2,
	Fn: func(ctx context.Context, args []core.Value) (bool, error) {
		return args[0].Type() == types.DateTime, nil
	},
}
