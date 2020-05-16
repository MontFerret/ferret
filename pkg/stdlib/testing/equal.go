package testing

import (
	"context"
	"fmt"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// Equal asserts equality of actual and expected values.
// @param (Mixed) - Actual value.
// @param (Mixed) - Expected value.
// @param (String) - Message to display on error.
var Equal = Assertion{
	DefaultMessage: func(args []core.Value) string {
		return fmt.Sprintf("be %s %s", EqualOp, formatValue(args[1]))
	},
	MinArgs: 2,
	MaxArgs: 3,
	Fn: func(ctx context.Context, args []core.Value) (bool, error) {
		return compare(args, EqualOp)
	},
}
