package testing

import (
	"context"
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/runtime"

	"github.com/MontFerret/ferret/v2/pkg/stdlib/testing/base"
)

// Tests whether the value is true.
// @param actual {Any} Value to test.
// @param message {String} Message to display on error.
// @return {Boolean} True when the configured assertion succeeds; otherwise an assertion error is returned.
var True = base.Assertion{
	DefaultMessage: func(args []runtime.Value) string {
		return fmt.Sprintf("be %s", base.FormatValue(runtime.True))
	},
	Args: base.Args{
		Min: 1,
		Max: 2,
	},
	Fn: func(ctx context.Context, args []runtime.Value) (bool, error) {
		return args[0] == runtime.True, nil
	},
}
