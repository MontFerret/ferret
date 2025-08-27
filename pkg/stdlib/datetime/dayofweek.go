package datetime

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// DATE_DAYOFWEEK returns number of the weekday from the date. Sunday is the 0th day of week.
// @param {DateTime} date - Source DateTime.
// @return {Int} - Number of the weekday.
func DateDayOfWeek(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, 1); err != nil {
		return runtime.None, err
	}

	if err := runtime.AssertDateTime(args[0]); err != nil {
		return runtime.None, err
	}

	wday := args[0].(runtime.DateTime).Weekday()

	return runtime.NewInt(int(wday)), nil
}
