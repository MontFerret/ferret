package testing

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// Equal asserts equality of actual and expected values.
// @param (Mixed) - Actual value.
// @param (Mixed) - Expected value.
// @param (String) - Message to display on error.
var Equal = Assertion{
	DefaultMessage: "",
	MessageArg:     0,
	MinArgs:        0,
	MaxArgs:        0,
	Fn: func(ctx context.Context, args []core.Value) (values.Boolean, error) {
		return compare(args, EqualOp)
	},
}
