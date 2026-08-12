package testing

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"

	"github.com/MontFerret/ferret/v2/pkg/stdlib/testing/base"
)

// fail returns an error.
// @param message {String} Message to display on error.
// @return {Boolean} No success value is produced because this assertion always fails.
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
