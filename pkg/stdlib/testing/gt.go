package testing

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// Gt asserts that an actual value is greater than an expected one.
// @param (Mixed) - Actual value.
// @param (Mixed) - Expected value.
// @param (String) - Message to display on error.
func Gt(_ context.Context, args ...core.Value) (core.Value, error) {
	return compare(args, GreaterOp)
}
