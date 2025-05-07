package testing

import (
	"context"
	"fmt"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/MontFerret/ferret/pkg/stdlib/collections"
	"github.com/MontFerret/ferret/pkg/stdlib/testing/base"
)

// INCLUDE asserts that haystack includes needle.
// @param {String | Array | Object | Iterable} actual - Haystack value.
// @param {Any} expected - Expected value.
// @param {String} [message] - Message to display on error.
var Include = base.Assertion{
	DefaultMessage: func(args []runtime.Value) string {
		return fmt.Sprintf("include %s", base.FormatValue(args[1]))
	},
	MinArgs: 2,
	MaxArgs: 3,
	Fn: func(ctx context.Context, args []runtime.Value) (bool, error) {
		haystack := args[0]
		needle := args[1]

		out, err := collections.Includes(ctx, haystack, needle)

		if err != nil {
			return false, err
		}

		return runtime.CompareValues(out, runtime.True) == 0, nil
	},
}
