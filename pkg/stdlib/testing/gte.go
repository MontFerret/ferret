package testing

import (
	"context"
	"fmt"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/stdlib/testing/base"
)

// Gte asserts that an actual value is greater than or equal to an expected one.
// @param actual {Any} - Actual value.
// @param expected {Any} - Expected value.
// @param message {String, optional} - Message to display on error.
var Gte = base.Assertion{
	DefaultMessage: func(args []core.Value) string {
		return fmt.Sprintf("be %s %s", base.GreaterOrEqualOp, base.FormatValue(args[1]))
	},
	MinArgs: 2,
	MaxArgs: 3,
	Fn: func(ctx context.Context, args []core.Value) (bool, error) {
		return base.GreaterOrEqualOp.Compare(args)
	},
}
