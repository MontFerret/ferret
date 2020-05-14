package testing

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/collections"
)

// Empty asserts that the target does not contain any values.
// @param value (Measurable|Binary|Object|Array|String) - Value to test.
// @param (String) - Message to display on error.
var Empty = Assertion{
	DefaultMessage: "empty",
	MessageArg:     1,
	MinArgs:        1,
	MaxArgs:        2,
	Fn: func(ctx context.Context, args []core.Value) (values.Boolean, error) {
		size, err := collections.Length(ctx, args[0])

		if err != nil {
			return values.False, err
		}

		return values.ToInt(size) == 0, nil
	},
}
