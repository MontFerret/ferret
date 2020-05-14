package testing

import (
	"context"
	"fmt"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// Gt asserts that an actual value is greater than an expected one.
// @param (Mixed) - Actual value.
// @param (Mixed) - Expected value.
// @param (String) - Message to display on error.
var Gt = Assertion{
	DefaultMessage: func(args []core.Value) string {
		return fmt.Sprintf("%s to be %s %s", args[0], GreaterOp, args[1])
	},
	MinArgs: 2,
	MaxArgs: 3,
	Fn: func(ctx context.Context, args []core.Value) (bool, error) {
		return compare(args, GreaterOp)
	},
}
