package testing

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"

	"github.com/MontFerret/ferret/v2/pkg/stdlib/testing/base"
)

// FAIL returns an error.
// @param {String} [message] - Message to display on error.
var Fail = base.Assertion{
	DefaultMessage: func(_ []runtime.Value) string {
		return "not fail"
	},
	Args: base.Args{
		Min: 0,
		Max: 1,
	},
	Fn: func(ctx context.Context, args []runtime.Value) (bool, error) {
		return false, nil
	},
}
