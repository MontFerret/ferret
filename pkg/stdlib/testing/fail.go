package testing

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/stdlib/testing/base"
)

// FAIL returns an error.
// @param {String} [message] - Message to display on error.
var Fail = base.Assertion{
	DefaultMessage: func(_ []core.Value) string {
		return "not fail"
	},
	MinArgs: 0,
	MaxArgs: 1,
	Fn: func(ctx context.Context, args []core.Value) (bool, error) {
		return false, nil
	},
}
