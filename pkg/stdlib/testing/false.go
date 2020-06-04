package testing

import (
	"context"
	"fmt"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// False asserts that value is false.
// @params actual (Mixed) - Value to test.
// @params message (String) - Message to display on error.
var False = Assertion{
	DefaultMessage: func(args []core.Value) string {
		return fmt.Sprintf("be %s", formatValue(values.False))
	},
	MinArgs: 1,
	MaxArgs: 2,
	Fn: func(ctx context.Context, args []core.Value) (bool, error) {
		return args[0] == values.False, nil
	},
}
