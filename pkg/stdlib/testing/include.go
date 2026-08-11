package testing

import (
	"context"
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/runtime"

	"github.com/MontFerret/ferret/v2/pkg/stdlib/collections"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/testing/base"
)

// Tests whether the actual container includes the expected value.
// @param actual {String | Array | Object | Iterable} Haystack value.
// @param expected {Any} Expected value.
// @param message {String} Message to display on error.
// @return {Boolean} True when the configured assertion succeeds; otherwise an assertion error is returned.
var Include = base.Assertion{
	DefaultMessage: func(args []runtime.Value) string {
		return fmt.Sprintf("include %s", base.FormatValue(args[1]))
	},
	Args: base.Args{
		Min: 2,
		Max: 3,
	},
	Fn: func(ctx context.Context, args []runtime.Value) (bool, error) {
		haystack := args[0]
		needle := args[1]

		out, err := collections.Includes(ctx, haystack, needle)

		if err != nil {
			return false, err
		}

		equal, err := runtime.EqualValues(ctx, out, runtime.True)
		return bool(equal), err
	},
}
