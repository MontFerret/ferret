package testing

import (
	"context"
	"fmt"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/testing/base"
)

// NONE asserts that value is none.
// @param {Any} actual - Value to test.
// @param {String} [message] - Message to display on error.
var None = base.Assertion{
	DefaultMessage: func(args []core.Value) string {
		return fmt.Sprintf("be %s", base.FormatValue(values.None))
	},
	MinArgs: 1,
	MaxArgs: 2,
	Fn: func(ctx context.Context, args []core.Value) (bool, error) {
		return args[0] == values.None, nil
	},
}
