package datetime

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// DATE_DAYOFWEEK returns number of the weekday from the date. Sunday is the 0th day of week.
// @param {DateTime} date - Source DateTime.
// @return {Int} - Number of the weekday.
func DateDayOfWeek(_ context.Context, args ...core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 1, 1); err != nil {
		return core.None, err
	}

	if err := core.AssertDateTime(args[0]); err != nil {
		return core.None, err
	}

	wday := args[0].(core.DateTime).Weekday()

	return core.NewInt(int(wday)), nil
}
