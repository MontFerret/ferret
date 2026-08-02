package types

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// ToDuration converts a value to a native duration.
// Numeric values are interpreted as milliseconds.
// @param {Any} value - A value coercible to Duration.
// @return {Duration} - A native duration value.
func ToDuration(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	return runtime.ToDuration(ctx, arg)
}
