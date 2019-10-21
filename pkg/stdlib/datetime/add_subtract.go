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

// DateAdd add amount given in unit to date.
// @params date (DateTime) - source date.
// @params amount (Int) - amount of units
// @params unit (String) - unit.
// @return (DateTime) - calculated date.
// The following units are available:
// * y, year, year
// * m, month, months
// * w, week, weeks
// * d, day, days
// * h, hour, hours
// * i, minute, minutes
// * s, second, seconds
// * f, millisecond, milliseconds
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

// DateSubtract subtract amount given in unit to date.
// @params date (DateTime) - source date.
// @params amount (Int) - amount of units
// @params unit (String) - unit.
// @return (DateTime) - calculated date.
// The following units are available:
// * y, year, year
// * m, month, months
// * w, week, weeks
// * d, day, days
// * h, hour, hours
// * i, minute, minutes
// * s, second, seconds
// * f, millisecond, milliseconds
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
