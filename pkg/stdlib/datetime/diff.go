package datetime

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// DATE_DIFF returns the difference between two dates in given time unit.
// @param {DateTime} date1 - First date.
// @param {DateTime} date2 - Second date.
// @param {String} unit - Time unit to return the difference in.
// @param {Boolean} [asFloat=False] - If true amount of unit will be as float.
// @return {Int | Float} - Difference between date1 and date2.
func DateDiff(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 3, 4)
	if err != nil {
		return values.None, err
	}

	err = core.ValidateValueTypePairs(
		core.NewPairValueType(args[0], types.DateTime),
		core.NewPairValueType(args[1], types.DateTime),
		core.NewPairValueType(args[2], types.String),
	)
	if err != nil {
		return values.None, err
	}

	date1 := args[0].(values.DateTime)
	date2 := args[1].(values.DateTime)
	unit := args[2].(values.String)
	isFloat := values.NewBoolean(false)

	if len(args) == 4 {
		err = core.ValidateType(args[3], types.Boolean)
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
	u, err := UnitFromString(unit)
	if err != nil {
		return -1, err
	}
	return nsec / u.Nanosecond(), nil
}
