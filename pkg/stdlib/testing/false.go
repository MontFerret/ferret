package testing

import (
	"context"
	"fmt"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/MontFerret/ferret/pkg/stdlib/testing/base"
)

// FALSE asserts that value is false.
// @param {Any}actual  - Second to test.
// @param {String} [message] - Message to display on error.
var False = base.Assertion{
	DefaultMessage: func(args []runtime.Value) string {
		return fmt.Sprintf("be %s", base.FormatValue(runtime.False))
	},
	MinArgs: 1,
	MaxArgs: 2,
	Fn: func(ctx context.Context, args []runtime.Value) (bool, error) {
		return args[0] == runtime.False, nil
	},
}
