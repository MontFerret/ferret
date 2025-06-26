package datetime

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// DATE_DIFF returns the difference between two dates in given time unit.
// @param {DateTime} date1 - First date.
// @param {DateTime} date2 - Second date.
// @param {String} unit - Time unit to return the difference in.
// @param {Boolean} [asFloat=False] - If true amount of unit will be as float.
// @return {Int | Float} - Difference between date1 and date2.
func DateDiff(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 3, 4); err != nil {
		return runtime.None, err
	}

	if err := runtime.AssertDateTime(args[0]); err != nil {
		return runtime.None, err
	}

	if err := runtime.AssertDateTime(args[1]); err != nil {
		return runtime.None, err
	}

	if err := runtime.AssertString(args[2]); err != nil {
		return runtime.None, err
	}

	date1 := args[0].(runtime.DateTime)
	date2 := args[1].(runtime.DateTime)
	unit := args[2].(runtime.String)
	isFloat := runtime.NewBoolean(false)

	if len(args) == 4 {
		if err := runtime.AssertBoolean(args[3]); err != nil {
			return runtime.None, err
		}

		isFloat = args[3].(runtime.Boolean)
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
