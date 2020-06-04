package testing

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/collections"
)

// Empty asserts that the target does not contain any values.
// @params actual (Measurable|Binary|Object|Array|String) - Value to test.
// @params message (String, optional) - Message to display on error.
var Empty = Assertion{
	DefaultMessage: func(_ []core.Value) string {
		return "be empty"
	},
	MinArgs: 1,
	MaxArgs: 2,
	Fn: func(ctx context.Context, args []core.Value) (bool, error) {
		size, err := collections.Length(ctx, args[0])

		if err != nil {
			return false, err
		}

		return values.ToInt(size) == 0, nil
	},
}
