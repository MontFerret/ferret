package datetime

import (
	"context"

	"github.com/pkg/errors"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// DATE_COMPARE checks if two partial dates match.
// @param {DateTime} date1 - First date.
// @param {DateTime} date2 - Second date.
// @param {String} unitRangeStart - Unit to start from.
// @param {String} [unitRangeEnd="millisecond"] - Unit to end with. Error will be returned if unitRangeStart unit less that unitRangeEnd.
// @return {Boolean} - True if the dates match, else false.
func DateCompare(_ context.Context, args ...core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 3, 4); err != nil {
		return core.None, err
	}

	if err := core.AssertDateTime(args[0]); err != nil {
		return core.None, err
	}

	if err := core.AssertDateTime(args[1]); err != nil {
		return core.None, err
	}

	if err := core.AssertString(args[2]); err != nil {
		return core.None, err
	}

	date1 := args[0].(core.DateTime)
	date2 := args[1].(core.DateTime)
	rangeStart := args[2].(core.String)
	rangeEnd := core.NewString("millisecond")

	if len(args) == 4 {
		if err := core.AssertString(args[3]); err != nil {
			return core.None, err
		}

		rangeEnd = args[3].(core.String)
	}

	unitStart, err := UnitFromString(rangeStart.String())
	if err != nil {
		return core.None, err
	}

	unitEnd, err := UnitFromString(rangeEnd.String())
	if err != nil {
		return core.None, err
	}

	if unitStart < unitEnd {
		return core.None, errors.Errorf("start unit less that end unit")
	}

	for u := unitEnd; u <= unitStart; u++ {
		if IsDatesEqual(date1.Time, date2.Time, u) {
			return core.NewBoolean(true), nil
		}
	}

	return core.NewBoolean(false), nil
}
