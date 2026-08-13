package datetime

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// date_diff returns the difference between two dates in given time unit.
// @param date1 {DateTime} First date.
// @param date2 {DateTime} Second date.
// @param unit {String} Time unit to return the difference in.
// @param asFloat {Boolean} If true amount of unit will be as float.
// @return {Int | Float} Difference between date1 and date2.
func DateDiff(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 3, 4); err != nil {
		return runtime.None, err
	}

	if len(args) == 3 {
		return dateDiff3(ctx, args[0], args[1], args[2])
	}

	return dateDiff4(ctx, args[0], args[1], args[2], args[3])
}

// date_diff returns the difference between two dates in given time unit.
// @param date1 {DateTime} First date.
// @param date2 {DateTime} Second date.
// @param unit {String} Time unit to return the difference in.
// @return {Int | Float} Difference between date1 and date2.
func dateDiff3(ctx context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
	return dateDiff4(ctx, arg1, arg2, arg3, runtime.False)
}

// date_diff returns the difference between two dates in given time unit.
// @param date1 {DateTime} First date.
// @param date2 {DateTime} Second date.
// @param unit {String} Time unit to return the difference in.
// @param asFloat {Boolean} If true amount of unit will be as float.
// @return {Int | Float} Difference between date1 and date2.
func dateDiff4(_ context.Context, arg1, arg2, arg3, arg4 runtime.Value) (runtime.Value, error) {
	date1, date2, unit, err := runtime.CastArgs3[runtime.DateTime, runtime.DateTime, runtime.String](arg1, arg2, arg3)

	if err != nil {
		return runtime.None, err
	}

	isFloat, err := runtime.CastArg[runtime.Boolean](arg4, 3)

	if err != nil {
		return runtime.None, err
	}

	if date1.Equal(date2.Time) {
		if isFloat {
			return runtime.NewFloat(0), nil
		}
		return runtime.NewInt(0), nil
	}

	var nsecDiff int64

	if date1.After(date2.Time) {
		nsecDiff = date1.Time.Sub(date2.Time).Nanoseconds()
	} else {
		nsecDiff = date2.Time.Sub(date1.Time).Nanoseconds()
	}

	unitDiff, err := nsecToUnit(float64(nsecDiff), unit.String())
	if err != nil {
		return runtime.None, err
	}

	if !isFloat {
		return runtime.NewInt(int(unitDiff)), nil
	}

	return runtime.NewFloat(unitDiff), nil
}

func nsecToUnit(nsec float64, unit string) (float64, error) {
	u, err := UnitFromString(unit)
	if err != nil {
		return -1, err
	}
	return nsec / u.Nanosecond(), nil
}
