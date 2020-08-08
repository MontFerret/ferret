package testing

import (
	"context"
	"fmt"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/stdlib/testing/base"
)

// LTE asserts that an actual value is lesser than or equal to an expected one.
// @param {Any} actual - Actual value.
// @param {Any} expected - Expected value.
// @param {String} [message] - Message to display on error.
var Lte = base.Assertion{
	DefaultMessage: func(args []core.Value) string {
		return fmt.Sprintf("be %s %s", base.LessOrEqualOp, base.FormatValue(args[1]))
	},
	MinArgs: 2,
	MaxArgs: 3,
	Fn: func(ctx context.Context, args []core.Value) (bool, error) {
		return base.LessOrEqualOp.Compare(args)
	},
}
