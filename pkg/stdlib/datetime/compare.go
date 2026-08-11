package datetime

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"

	"errors"
)

// DATE_COMPARE checks if two partial dates match.
// @param date1 {DateTime} First date.
// @param date2 {DateTime} Second date.
// @param unitRangeStart {String} Unit to start from.
// @param unitRangeEnd {String} Unit to end with. Error will be returned if unitRangeStart unit less that unitRangeEnd.
// @return {Boolean} True if the dates match, else false.
func DateCompare(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 3, 4); err != nil {
		return runtime.None, err
	}

	if len(args) == 3 {
		return dateCompare3(ctx, args[0], args[1], args[2])
	}

	return dateCompare4(ctx, args[0], args[1], args[2], args[3])
}

// DATE_COMPARE checks if two partial dates match.
// @param date1 {DateTime} First date.
// @param date2 {DateTime} Second date.
// @param unitRangeStart {String} Unit to start from.
// @return {Boolean} True if the dates match, else false.
func dateCompare3(ctx context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
	return dateCompare4(ctx, arg1, arg2, arg3, runtime.NewString("millisecond"))
}

// DATE_COMPARE checks if two partial dates match.
// @param date1 {DateTime} First date.
// @param date2 {DateTime} Second date.
// @param unitRangeStart {String} Unit to start from.
// @param unitRangeEnd {String} Unit to end with. Error will be returned if unitRangeStart unit less that unitRangeEnd.
// @return {Boolean} True if the dates match, else false.
func dateCompare4(_ context.Context, arg1, arg2, arg3, arg4 runtime.Value) (runtime.Value, error) {
	if err := runtime.AssertDateTime(arg1); err != nil {
		return runtime.None, err
	}

	if err := runtime.AssertDateTime(arg2); err != nil {
		return runtime.None, err
	}

	if err := runtime.AssertString(arg3); err != nil {
		return runtime.None, err
	}

	if err := runtime.AssertString(arg4); err != nil {
		return runtime.None, err
	}

	date1 := arg1.(runtime.DateTime)
	date2 := arg2.(runtime.DateTime)
	rangeStart := arg3.(runtime.String)
	rangeEnd := arg4.(runtime.String)

	unitStart, err := UnitFromString(rangeStart.String())
	if err != nil {
		return runtime.None, err
	}

	unitEnd, err := UnitFromString(rangeEnd.String())
	if err != nil {
		return runtime.None, err
	}

	if unitStart < unitEnd {
		return runtime.None, errors.New("start unit less that end unit")
	}

	for u := unitEnd; u <= unitStart; u++ {
		if IsDatesEqual(date1.Time, date2.Time, u) {
			return runtime.NewBoolean(true), nil
		}
	}

	return runtime.NewBoolean(false), nil
}
