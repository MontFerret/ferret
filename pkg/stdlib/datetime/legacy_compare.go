package datetime

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// date_compare requires every component from start through millisecond to match.
// The component order is year, month, week, day, hour, minute, second, millisecond.
// Each input uses its own calendar; week includes the ISO week-year.
// @param left {DateTime} First date.
// @param right {DateTime} Second date.
// @param start {String} Coarsest component in the inclusive range.
// @return {Boolean} Whether all selected components match.
// @throws {InvalidArgument} The unit is unknown.
// @deprecated Use datetime::same for precision equality; it has different semantics from component ranges.
func legacyCompare3(ctx context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
	return legacyCompare4(ctx, arg1, arg2, arg3, runtime.String("millisecond"))
}

// date_compare requires every component in the inclusive start/end range to match.
// The component order is year, month, week, day, hour, minute, second, millisecond.
// Each input uses its own calendar; week includes the ISO week-year.
// @param left {DateTime} First date.
// @param right {DateTime} Second date.
// @param start {String} Coarsest component in the inclusive range.
// @param end {String} Finest component; must not be coarser than start.
// @return {Boolean} Whether all selected components match.
// @throws {InvalidArgument} A unit is unknown or the range is reversed.
// @deprecated Use datetime::same for precision equality; it has different semantics from component ranges.
func legacyCompare4(_ context.Context, arg1, arg2, arg3, arg4 runtime.Value) (runtime.Value, error) {
	left, right, start, err := datePairUnit(arg1, arg2, arg3)
	if err != nil {
		return runtime.None, err
	}

	end, err := parseUnit(arg4, 3)
	if err != nil {
		return runtime.None, err
	}

	if start < end {
		return runtime.None, runtime.ArgError(runtime.Error(runtime.ErrInvalidArgument, "end unit must not be coarser than start unit"), 3)
	}

	for current := end; current <= start; current++ {
		if !sameComponent(left.Time, right.Time, current) {
			return runtime.False, nil
		}
	}

	return runtime.True, nil
}
