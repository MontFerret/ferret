package types

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// to_datetime converts one native DateTime or RFC3339 value.
// @param value {DateTime|String} A DateTime or RFC3339 string.
// @return {DateTime} Parsed DateTime.
func ToDateTime(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	return runtime.ToDateTime(ctx, arg)
}

// to_datetime returns a native DateTime, parses an RFC3339 string, or converts
// a numeric Unix epoch value when an explicit unit is provided.
// @param value {DateTime|String|Int|Float} A DateTime, RFC3339 string, or numeric Unix epoch value.
// @param unit {String} Epoch unit: s, ms, us, or ns. Aliases include sec, second, seconds, millisecond, milliseconds, µs, μs, microsecond, microseconds, nanosecond, and nanoseconds. Valid only for numeric values.
// @return {DateTime} Parsed DateTime.
func toDateTime2(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	return runtime.ToDateTimeEpoch(ctx, arg1, arg2)
}
