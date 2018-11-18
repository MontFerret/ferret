package datetime

import (
	"context"
	"strings"
	"time"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/pkg/errors"
)

var (
	sliceDateTime   = []core.Type{core.DateTimeType}
	sliceIntType    = []core.Type{core.IntType}
	sliceStringType = []core.Type{core.StringType}

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

	dt, err := addUnit(date, int(amount), unit.String())
	if err != nil {
		return values.None, err
	}

	return dt, nil
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

	dt, err := addUnit(date, -1*int(amount), unit.String())
	if err != nil {
		return values.None, err
	}

	return dt, nil
}

func getArgs(args []core.Value) (values.DateTime, values.Int, values.String, error) {
	err := core.ValidateArgs(args, 3, 3)
	if err != nil {
		return emptyDateTime, emptyInt, emptyString, err
	}

	err = core.ValidateValueTypePairs(
		core.PairValueType{Value: args[0], Types: sliceDateTime},
		core.PairValueType{Value: args[1], Types: sliceIntType},
		core.PairValueType{Value: args[2], Types: sliceStringType},
	)
	if err != nil {
		return emptyDateTime, emptyInt, emptyString, err
	}

	date := args[0].(values.DateTime)
	amount := args[1].(values.Int)
	unit := args[2].(values.String)

	return date, amount, unit, nil
}

func addUnit(dt values.DateTime, amount int, unit string) (values.DateTime, error) {
	switch strings.ToLower(unit) {
	case "y", "year", "years":
		return values.NewDateTime(dt.AddDate(amount*1, 0, 0)), nil
	case "m", "month", "months":
		return values.NewDateTime(dt.AddDate(0, amount*1, 0)), nil
	case "w", "week", "weeks":
		return values.NewDateTime(dt.AddDate(0, 0, amount*7)), nil
	case "d", "day", "days":
		return values.NewDateTime(dt.AddDate(0, 0, amount*1)), nil
	case "h", "hour", "hours":
		return values.NewDateTime(dt.Add(time.Duration(amount) * time.Hour)), nil
	case "i", "minute", "minutes":
		return values.NewDateTime(dt.Add(time.Duration(amount) * time.Minute)), nil
	case "s", "second", "seconds":
		return values.NewDateTime(dt.Add(time.Duration(amount) * time.Second)), nil
	case "f", "millisecond", "milliseconds":
		return values.NewDateTime(dt.Add(time.Duration(amount) * time.Millisecond)), nil
	}
	return values.DateTime{}, errors.Errorf("no such unit '%s'", unit)
}
