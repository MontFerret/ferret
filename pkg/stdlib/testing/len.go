package testing

import (
	"context"
	"fmt"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/stdlib/collections"
	"github.com/MontFerret/ferret/pkg/stdlib/testing/base"
)

// LEN asserts that a measurable value has a length or size with the expected value.
// @param {Measurable} actual - Measurable value.
// @param {Int} length - Target length.
// @param {String} [message] - Message to display on error.
var Len = base.Assertion{
	DefaultMessage: func(args []core.Value) string {
		return fmt.Sprintf("has size %s", args[1])
	},
	MinArgs: 2,
	MaxArgs: 3,
	Fn: func(ctx context.Context, args []core.Value) (bool, error) {
		col := args[0]
		size := args[1]

		out, err := collections.Length(ctx, col)

		if err != nil {
			return false, err
		}

		return out.Compare(size) == 0, nil
	},
}
