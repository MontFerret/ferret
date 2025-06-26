package datetime

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/pkg/errors"
)

// DATE_COMPARE checks if two partial dates match.
// @param {DateTime} date1 - First date.
// @param {DateTime} date2 - Second date.
// @param {String} unitRangeStart - Unit to start from.
// @param {String} [unitRangeEnd="millisecond"] - Unit to end with. Error will be returned if unitRangeStart unit less that unitRangeEnd.
// @return {Boolean} - True if the dates match, else false.
func DateCompare(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 3, 4); err != nil {
		return runtime.None, err
	}

	if err := runtime.AssertDateTime(args[0]); err != nil {
		return runtime.None, err
	}

	if err := runtime.AssertDateTime(args[1]); err != nil {
		return runtime.None, err
	}

	if err := runtime.AssertString(args[2]); err != nil {
		return runtime.None, err
	}

	date1 := args[0].(runtime.DateTime)
	date2 := args[1].(runtime.DateTime)
	rangeStart := args[2].(runtime.String)
	rangeEnd := runtime.NewString("millisecond")

	if len(args) == 4 {
		if err := runtime.AssertString(args[3]); err != nil {
			return runtime.None, err
		}

		rangeEnd = args[3].(runtime.String)
	}

	unitStart, err := UnitFromString(rangeStart.String())
	if err != nil {
		return runtime.None, err
	}

	unitEnd, err := UnitFromString(rangeEnd.String())
	if err != nil {
		return runtime.None, err
	}

	if unitStart < unitEnd {
		return runtime.None, errors.Errorf("start unit less that end unit")
	}

	for u := unitEnd; u <= unitStart; u++ {
		if IsDatesEqual(date1.Time, date2.Time, u) {
			return runtime.NewBoolean(true), nil
		}
	}

	return runtime.NewBoolean(false), nil
}
