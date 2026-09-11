package datetime

import (
	"context"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// Diff returns the signed elapsed difference right - left in fixed units.
// It uses checked runtime subtraction, limited to Duration's roughly 292-year range.
// @param left {DateTime} Starting instant.
// @param right {DateTime} Ending instant.
// @param unit {String} Millisecond, second, minute, or hour; plurals are accepted. Calendar units are rejected.
// @return {Float} Signed fractional elapsed units; reversing arguments reverses the sign.
// @throws {TypeError} An argument has an invalid type.
// @throws {InvalidArgument} The unit is unknown or is a calendar unit.
// @throws {RangeError} The elapsed interval cannot fit a runtime Duration.
func Diff(ctx context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
	elapsed, divisor, err := elapsedDifference(ctx, arg1, arg2, arg3)
	if err != nil {
		return runtime.None, err
	}

	return runtime.Float(float64(elapsed) / float64(divisor)), nil
}

func elapsedDifference(ctx context.Context, arg1, arg2, arg3 runtime.Value) (time.Duration, time.Duration, error) {
	_, _, u, err := datePairUnit(arg1, arg2, arg3)
	if err != nil {
		return 0, 0, err
	}

	divisor, ok := fixedDuration(u)
	if !ok {
		return 0, 0, runtime.ArgError(runtime.Error(runtime.ErrInvalidArgument, "datetime diff supports only milliseconds, seconds, minutes, and hours; calendar units are not supported"), 2)
	}

	// Reuse the validated values without boxing the DateTimes again.
	result, err := runtime.Subtract(ctx, arg2, arg1)
	if err != nil {
		return 0, 0, err
	}

	return time.Duration(result.(runtime.Duration)), divisor, nil
}
