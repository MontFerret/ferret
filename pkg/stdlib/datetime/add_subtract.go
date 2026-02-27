package datetime

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

var (
	emptyDateTime runtime.DateTime
	emptyInt      runtime.Int
	emptyString   runtime.String
)

// DATE_ADD adds amount given in unit to date.
// The following units are available:
// * y, year, year
// * m, month, months
// * w, week, weeks
// * d, day, days
// * h, hour, hours
// * i, minute, minutes
// * s, second, seconds
// * f, millisecond, milliseconds
// @param {DateTime} date - Source date.
// @param {Int} amount - Amount of units
// @param {String} unit - Unit.
// @return {DateTime} - Calculated date.
func DateAdd(_ context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
	date, amount, unit, err := castArgs(arg1, arg2, arg3)
	if err != nil {
		return runtime.None, err
	}

	u, err := UnitFromString(unit.String())
	if err != nil {
		return runtime.None, err
	}

	tm := AddUnit(date.Time, int(amount), u)

	return runtime.NewDateTime(tm), nil
}

// DATE_SUBTRACT subtract amount given in unit to date.
// The following units are available:
// * y, year, year
// * m, month, months
// * w, week, weeks
// * d, day, days
// * h, hour, hours
// * i, minute, minutes
// * s, second, seconds
// * f, millisecond, milliseconds
// @param {DateTime} date - source date.
// @param {Int} amount - amount of units
// @param {String} unit - unit.
// @return {DateTime} - calculated date.
func DateSubtract(_ context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
	date, amount, unit, err := castArgs(arg1, arg2, arg3)
	if err != nil {
		return runtime.None, err
	}

	u, err := UnitFromString(unit.String())
	if err != nil {
		return runtime.None, err
	}

	tm := AddUnit(date.Time, -1*int(amount), u)

	return runtime.NewDateTime(tm), nil
}

func castArgs(arg1, arg2, arg3 runtime.Value) (runtime.DateTime, runtime.Int, runtime.String, error) {
	return runtime.CastArgs3[runtime.DateTime, runtime.Int, runtime.String](arg1, arg2, arg3)
}
