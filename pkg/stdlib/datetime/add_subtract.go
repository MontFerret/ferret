package datetime

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
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
func DateAdd(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	date, amount, unit, err := getArgs(args)
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
func DateSubtract(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	date, amount, unit, err := getArgs(args)
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

func getArgs(args []runtime.Value) (runtime.DateTime, runtime.Int, runtime.String, error) {
	if err := runtime.ValidateArgs(args, 3, 3); err != nil {
		return emptyDateTime, emptyInt, emptyString, err
	}

	if err := runtime.AssertDateTime(args[0]); err != nil {
		return emptyDateTime, emptyInt, emptyString, err
	}

	if err := runtime.AssertInt(args[1]); err != nil {
		return emptyDateTime, emptyInt, emptyString, err
	}

	if err := runtime.AssertString(args[2]); err != nil {
		return emptyDateTime, emptyInt, emptyString, err
	}

	date := args[0].(runtime.DateTime)
	amount := args[1].(runtime.Int)
	unit := args[2].(runtime.String)

	return date, amount, unit, nil
}
