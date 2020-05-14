package testing

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// Fail returns an error.
// @param (String) - Message to display on error.
var Fail = Assertion{
	DefaultMessage: func(_ []core.Value) string {
		return "not fail"
	},
	MinArgs: 0,
	MaxArgs: 1,
	Fn: func(ctx context.Context, args []core.Value) (bool, error) {
		return false, nil
	},
}
