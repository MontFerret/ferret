package types

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// ToDuration converts a value using the language's broad Duration coercion rules.
// Numeric values are interpreted as milliseconds; Duration strings, NONE,
// Booleans, and supported list shapes are also accepted.
// @param {Any} value - A value coercible to Duration.
// @return {Duration} - A native duration value.
func ToDuration(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	return runtime.ToDuration(ctx, arg)
}
