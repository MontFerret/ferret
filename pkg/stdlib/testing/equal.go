package testing

import (
	"context"
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/runtime"

	"github.com/MontFerret/ferret/v2/pkg/stdlib/testing/base"
)

// EQUAL asserts equality of actual and expected values.
// @param {Any} actual - Actual value.
// @param {Any} expected - Expected value.
// @param {String} [message] - Message to display on error.
var Equal = base.Assertion{
	DefaultMessage: func(args []runtime.Value) string {
		return fmt.Sprintf("be %s %s", base.EqualOp, base.FormatValue(args[1]))
	},
	MinArgs: 2,
	MaxArgs: 3,
	Fn: func(ctx context.Context, args []runtime.Value) (bool, error) {
		return base.EqualOp.Compare(args)
	},
}
