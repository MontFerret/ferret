package datetime

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

var (
	emptyDateTime values.DateTime
	emptyInt      values.Int
	emptyString   values.String
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
		return values.None, err
	}

	u, err := UnitFromString(unit.String())
	if err != nil {
		return values.None, err
	}

	tm := AddUnit(date.Time, int(amount), u)

	return values.NewDateTime(tm), nil
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
		return values.None, err
	}

	u, err := UnitFromString(unit.String())
	if err != nil {
		return values.None, err
	}

	tm := AddUnit(date.Time, -1*int(amount), u)

	return values.NewDateTime(tm), nil
}

func getArgs(args []core.Value) (values.DateTime, values.Int, values.String, error) {
	if err := core.ValidateArgs(args, 3, 3); err != nil {
		return emptyDateTime, emptyInt, emptyString, err
	}

	if err := values.AssertDateTime(args[0]); err != nil {
		return emptyDateTime, emptyInt, emptyString, err
	}

	if err := values.AssertInt(args[1]); err != nil {
		return emptyDateTime, emptyInt, emptyString, err
	}

	if err := values.AssertString(args[2]); err != nil {
		return emptyDateTime, emptyInt, emptyString, err
	}

	date := args[0].(values.DateTime)
	amount := args[1].(values.Int)
	unit := args[2].(values.String)

	return date, amount, unit, nil
}
