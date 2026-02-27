package datetime

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// DATE_DAYOFWEEK returns number of the weekday from the date. Sunday is the 0th day of week.
// @param {DateTime} date - Source DateTime.
// @return {Int} - Number of the weekday.
func DateDayOfWeek(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	dt, err := runtime.CastArg[runtime.DateTime](arg, 0)

	if err != nil {
		return runtime.None, err
	}

	wday := dt.Weekday()

	return runtime.Int(wday), nil
}
