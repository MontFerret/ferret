package datetime

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
	"github.com/pkg/errors"
)

// DateCompare check if two partial dates match.
// @params date1, date2 (DateTime) - comparable dates.
// @params unitRangeStart (String) - unit to start from.
// @params unitRangeEnd (String, Optional) - unit to end with.
// Error will be returned if unitRangeStart unit less that unitRangeEnd.
// @return (Boolean) - true if the dates match, else false.
func DateCompare(_ context.Context, args ...core.Value) (core.Value, error) {
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
	rangeStart := args[2].(values.String)
	rangeEnd := values.NewString("millisecond")

	if len(args) == 4 {
		err = core.ValidateType(args[3], types.String)

		if err != nil {
			return values.None, err
		}

		rangeEnd = args[3].(values.String)
	}

	unitStart, err := UnitFromString(rangeStart.String())
	if err != nil {
		return values.None, err
	}

	unitEnd, err := UnitFromString(rangeEnd.String())
	if err != nil {
		return values.None, err
	}

	if unitStart < unitEnd {
		return values.None, errors.Errorf("start unit less that end unit")
	}

	for u := unitEnd; u <= unitStart; u++ {
		if IsDatesEqual(date1.Time, date2.Time, u) {
			return values.NewBoolean(true), nil
		}
	}

	return values.NewBoolean(false), nil
}
