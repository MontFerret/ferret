package datetime

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// date_diff returns signed whole elapsed units for right - left, truncated toward zero.
// @param left {DateTime} Starting instant.
// @param right {DateTime} Ending instant.
// @param unit {String} Millisecond, second, minute, or hour; plurals are accepted. Calendar units are rejected.
// @return {Int} Signed whole elapsed units.
// @throws {InvalidArgument} The unit is unknown or is a calendar unit.
// @throws {RangeError} The elapsed interval cannot fit a runtime Duration.
// @deprecated Use datetime::diff instead; it always returns a Float.
func legacyDiff3(ctx context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
	return legacyDiff4(ctx, arg1, arg2, arg3, runtime.False)
}

// date_diff returns signed elapsed units for right - left, optionally fractional.
// @param left {DateTime} Starting instant.
// @param right {DateTime} Ending instant.
// @param unit {String} Millisecond, second, minute, or hour; plurals are accepted. Calendar units are rejected.
// @param asFloat {Boolean} True returns fractional units; false truncates toward zero.
// @return {Int | Float} Signed elapsed units.
// @throws {InvalidArgument} The unit is unknown or is a calendar unit.
// @throws {RangeError} The elapsed interval cannot fit a runtime Duration.
// @deprecated Use datetime::diff instead; it always returns a Float without an asFloat flag.
func legacyDiff4(ctx context.Context, arg1, arg2, arg3, arg4 runtime.Value) (runtime.Value, error) {
	asFloat, err := runtime.CastArg[runtime.Boolean](arg4, 3)
	if err != nil {
		return runtime.None, err
	}

	elapsed, divisor, err := elapsedDifference(ctx, arg1, arg2, arg3)
	if err != nil {
		return runtime.None, err
	}

	if asFloat {
		return runtime.Float(float64(elapsed) / float64(divisor)), nil
	}

	return runtime.Int(elapsed / divisor), nil
}
