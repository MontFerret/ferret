package testing

import (
	"context"
	"fmt"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/stdlib/testing/base"
)

// Lt asserts that an actual value is lesser than an expected one.
// @param actual (Value) - Actual value.
// @param expected (Value) - Expected value.
// @param message (String, optional) - Message to display on error.
var Lt = base.Assertion{
	DefaultMessage: func(args []core.Value) string {
		return fmt.Sprintf("be %s %s", base.LessOp, base.FormatValue(args[1]))
	},
	MinArgs: 2,
	MaxArgs: 3,
	Fn: func(ctx context.Context, args []core.Value) (bool, error) {
		return base.LessOp.Compare(args)
	},
}
