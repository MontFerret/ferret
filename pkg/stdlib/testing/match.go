package testing

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/strings"
)

// Match asserts that value matches the regular expression.
// @param (Mixed) - Actual value.
// @param (String) - Regular expression.
// @param (String) - Message to display on error.
var Match = Assertion{
	DefaultMessage: func(args []core.Value) string {
		return "match regular expression"
	},
	MinArgs: 2,
	MaxArgs: 3,
	Fn: func(ctx context.Context, args []core.Value) (bool, error) {
		value := args[0]
		regexp := args[1]

		out, err := strings.RegexMatch(ctx, value, regexp)

		if err != nil {
			return false, err
		}

		return out.Compare(values.True) == 0, nil
	},
}
