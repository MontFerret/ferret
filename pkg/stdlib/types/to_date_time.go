package types

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// TO_DATETIME returns a native DateTime or parses an RFC3339 string.
// @param {DateTime|String} value - A native DateTime or RFC3339 string.
// @return {DateTime} - Parsed DateTime.
func ToDateTime(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	return runtime.ToDateTime(ctx, arg)
}
