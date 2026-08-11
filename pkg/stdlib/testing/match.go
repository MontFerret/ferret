package testing

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"

	"github.com/MontFerret/ferret/v2/pkg/stdlib/strings"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/testing/base"
)

// Tests whether the value matches the regular expression.
// @param actual {Any} Actual value.
// @param expression {String} Regular expression.
// @param message {String} Message to display on error.
// @return {Boolean} True when the configured assertion succeeds; otherwise an assertion error is returned.
var Match = base.Assertion{
	DefaultMessage: func(args []runtime.Value) string {
		return "match regular expression"
	},
	Args: base.Args{
		Min: 2,
		Max: 3,
	},
	Fn: func(ctx context.Context, args []runtime.Value) (bool, error) {
		value := args[0]
		regexp := args[1]

		out, err := strings.RegexTest(ctx, value, regexp)

		if err != nil {
			return false, err
		}

		equal, err := runtime.EqualValues(ctx, out, runtime.True)
		return bool(equal), err
	},
}
