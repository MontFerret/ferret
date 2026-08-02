package types

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// ToDateTime converts one native DateTime or RFC3339 value.
// Numeric epoch conversion is exposed to FQL through the registered variadic adapter.
func ToDateTime(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	return runtime.ToDateTime(ctx, arg)
}

// TO_DATETIME returns a native DateTime, parses an RFC3339 string, or converts
// a numeric Unix epoch value when an explicit unit is provided.
// @param {DateTime|String|Int|Float} value - A DateTime, RFC3339 string, or numeric Unix epoch value.
// @param {String} [unit] - Epoch unit: s, ms, us, or ns. Aliases include sec, second, seconds, millisecond, milliseconds, µs, μs, microsecond, microseconds, nanosecond, and nanoseconds. Valid only for numeric values.
// @return {DateTime} - Parsed DateTime.
func toDateTime(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, 2); err != nil {
		return runtime.None, err
	}

	if len(args) == 1 {
		return runtime.ToDateTime(ctx, args[0])
	}

	return runtime.ToDateTimeEpoch(ctx, args[0], args[1])
}
