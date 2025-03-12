package datetime

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

var (
	emptyDateTime core.DateTime
	emptyInt      core.Int
	emptyString   core.String
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
func DateAdd(_ context.Context, args ...core.Value) (core.Value, error) {
	date, amount, unit, err := getArgs(args)
	if err != nil {
		return core.None, err
	}

	u, err := UnitFromString(unit.String())
	if err != nil {
		return core.None, err
	}

	tm := AddUnit(date.Time, int(amount), u)

	return core.NewDateTime(tm), nil
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
func DateSubtract(_ context.Context, args ...core.Value) (core.Value, error) {
	date, amount, unit, err := getArgs(args)
	if err != nil {
		return core.None, err
	}

	u, err := UnitFromString(unit.String())
	if err != nil {
		return core.None, err
	}

	tm := AddUnit(date.Time, -1*int(amount), u)

	return core.NewDateTime(tm), nil
}

func getArgs(args []core.Value) (core.DateTime, core.Int, core.String, error) {
	if err := core.ValidateArgs(args, 3, 3); err != nil {
		return emptyDateTime, emptyInt, emptyString, err
	}

	if err := core.AssertDateTime(args[0]); err != nil {
		return emptyDateTime, emptyInt, emptyString, err
	}

	if err := core.AssertInt(args[1]); err != nil {
		return emptyDateTime, emptyInt, emptyString, err
	}

	if err := core.AssertString(args[2]); err != nil {
		return emptyDateTime, emptyInt, emptyString, err
	}

	date := args[0].(core.DateTime)
	amount := args[1].(core.Int)
	unit := args[2].(core.String)

	return date, amount, unit, nil
}
