package testing

import (
	"context"
	"fmt"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// Gte asserts that an actual value is lesser than or equal to an expected one.
// @param (Mixed) - Actual value.
// @param (Mixed) - Expected value.
// @param (String) - Message to display on error.
var Lte = Assertion{
	DefaultMessage: func(args []core.Value) string {
		return fmt.Sprintf("be %s %s", LessOrEqualOp, formatValue(args[1]))
	},
	MinArgs: 2,
	MaxArgs: 3,
	Fn: func(ctx context.Context, args []core.Value) (bool, error) {
		return compare(args, LessOrEqualOp)
	},
}
