package datetime

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
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
	err := core.ValidateArgs(args, 3, 3)
	if err != nil {
		return emptyDateTime, emptyInt, emptyString, err
	}

	err = core.ValidateValueTypePairs(
		core.NewPairValueType(args[0], types.DateTime),
		core.NewPairValueType(args[1], types.Int),
		core.NewPairValueType(args[2], types.String),
	)
	if err != nil {
		return emptyDateTime, emptyInt, emptyString, err
	}

	date := args[0].(values.DateTime)
	amount := args[1].(values.Int)
	unit := args[2].(values.String)

	return date, amount, unit, nil
}
