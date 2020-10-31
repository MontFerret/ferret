package testing

import (
	"context"
	"fmt"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/testing/base"
)

// FALSE asserts that value is false.
// @param {Any}actual  - Value to test.
// @param {String} [message] - Message to display on error.
var False = base.Assertion{
	DefaultMessage: func(args []core.Value) string {
		return fmt.Sprintf("be %s", base.FormatValue(values.False))
	},
	MinArgs: 1,
	MaxArgs: 2,
	Fn: func(ctx context.Context, args []core.Value) (bool, error) {
		return args[0] == values.False, nil
	},
}
