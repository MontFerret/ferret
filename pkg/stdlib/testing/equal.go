package testing

import (
	"context"
	"fmt"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/stdlib/testing/base"
)

// Equal asserts equality of actual and expected values.
// @param actual {Any} - Actual value.
// @param expected {Any} - Expected value.
// @param message {String, optional} - Message to display on error.
var Equal = base.Assertion{
	DefaultMessage: func(args []core.Value) string {
		return fmt.Sprintf("be %s %s", base.EqualOp, base.FormatValue(args[1]))
	},
	MinArgs: 2,
	MaxArgs: 3,
	Fn: func(ctx context.Context, args []core.Value) (bool, error) {
		return base.EqualOp.Compare(args)
	},
}
