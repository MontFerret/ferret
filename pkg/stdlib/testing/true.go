package testing

import (
	"context"
	"fmt"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// True asserts that value is true.
// @params actual (Mixed) - Value to test.
// @params message (String, optional) - Message to display on error.
var True = Assertion{
	DefaultMessage: func(args []core.Value) string {
		return fmt.Sprintf("be %s", formatValue(values.True))
	},
	MinArgs: 1,
	MaxArgs: 2,
	Fn: func(ctx context.Context, args []core.Value) (bool, error) {
		return args[0] == values.True, nil
	},
}
