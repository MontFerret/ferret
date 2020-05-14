package testing

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// NotEqual asserts inequality of actual and expected values.
// @param (Mixed) - Actual value.
// @param (Mixed) - Expected value.
// @param (String) - Message to display on error.
func NotEqual(_ context.Context, args ...core.Value) (core.Value, error) {
	return compare(args, NotEqualOp)
}
