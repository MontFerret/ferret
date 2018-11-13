package datetime

import (
	"context"
	"strings"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/pkg/errors"
)

// DateDiff returns the difference between two dates in given time unit.
// @params date1   (DateTime) - first DateTime.
// @params date2   (DateTime) - second DateTime.
// @params unit    (String)   - time unit to return the difference in.
// @params asFloat (Boolean, optional) - if true amount of unit will be as float.
// @return (Int, Float) - difference between date1 and date2.
func DateDiff(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 3, 4)
	if err != nil {
		return values.None, err
	}

	err = core.ValidateValueTypePairs(
		core.PairValueType{Value: args[0], Types: sliceDateTime},
		core.PairValueType{Value: args[1], Types: sliceDateTime},
		core.PairValueType{Value: args[2], Types: sliceStringType},
	)
	if err != nil {
		return values.None, err
	}

	date1 := args[0].(values.DateTime)
	date2 := args[1].(values.DateTime)
	unit := args[2].(values.String)
	isFloat := values.NewBoolean(false)

	if len(args) == 4 {
		err = core.ValidateType(args[3], core.BooleanType)
		if err != nil {
			return values.None, err
		}
		isFloat = args[3].(values.Boolean)
	}

	if date1.Equal(date2.Time) {
		if isFloat {
			return values.NewFloat(0), nil
		}
		return values.NewInt(0), nil
	}

	var nsecDiff int64

	if date1.After(date2.Time) {
		nsecDiff = date1.Time.Sub(date2.Time).Nanoseconds()
	} else {
		nsecDiff = date2.Time.Sub(date1.Time).Nanoseconds()
	}

	unitDiff, err := nsecToUnit(float64(nsecDiff), unit.String())
	if err != nil {
		return values.None, err
	}

	if !isFloat {
		return values.NewInt(int(unitDiff)), nil
	}

	return values.NewFloat(unitDiff), nil
}

func nsecToUnit(nsec float64, unit string) (float64, error) {
	switch strings.ToLower(unit) {
	case "y", "year", "years":
		return nsec / 31536e12, nil
	case "m", "month", "months":
		return nsec / 26784e11, nil
	case "w", "week", "weeks":
		return nsec / 6048e11, nil
	case "d", "day", "days":
		return nsec / 864e11, nil
	case "h", "hour", "hours":
		return nsec / 36e11, nil
	case "i", "minute", "minutes":
		return nsec / 6e10, nil
	case "s", "second", "seconds":
		return nsec / 1e9, nil
	case "f", "millisecond", "milliseconds":
		return nsec / 1e6, nil
	}
	return -1, errors.Errorf("no such unit '%s'", unit)
}
