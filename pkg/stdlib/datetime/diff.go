package datetime

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// DATE_DIFF returns the difference between two dates in given time unit.
// @param {DateTime} date1 - First date.
// @param {DateTime} date2 - Second date.
// @param {String} unit - Time unit to return the difference in.
// @param {Boolean} [asFloat=False] - If true amount of unit will be as float.
// @return {Int | Float} - Difference between date1 and date2.
func DateDiff(_ context.Context, args ...core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 3, 4); err != nil {
		return core.None, err
	}

	if err := core.AssertDateTime(args[0]); err != nil {
		return core.None, err
	}

	if err := core.AssertDateTime(args[1]); err != nil {
		return core.None, err
	}

	if err := core.AssertString(args[2]); err != nil {
		return core.None, err
	}

	date1 := args[0].(core.DateTime)
	date2 := args[1].(core.DateTime)
	unit := args[2].(core.String)
	isFloat := core.NewBoolean(false)

	if len(args) == 4 {
		if err := core.AssertBoolean(args[3]); err != nil {
			return core.None, err
		}

		isFloat = args[3].(core.Boolean)
	}

	if date1.Equal(date2.Time) {
		if isFloat {
			return core.NewFloat(0), nil
		}
		return core.NewInt(0), nil
	}

	var nsecDiff int64

	if date1.After(date2.Time) {
		nsecDiff = date1.Time.Sub(date2.Time).Nanoseconds()
	} else {
		nsecDiff = date2.Time.Sub(date1.Time).Nanoseconds()
	}

	unitDiff, err := nsecToUnit(float64(nsecDiff), unit.String())
	if err != nil {
		return core.None, err
	}

	if !isFloat {
		return core.NewInt(int(unitDiff)), nil
	}

	return core.NewFloat(unitDiff), nil
}

func nsecToUnit(nsec float64, unit string) (float64, error) {
	u, err := UnitFromString(unit)
	if err != nil {
		return -1, err
	}
	return nsec / u.Nanosecond(), nil
}
