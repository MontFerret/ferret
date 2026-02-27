package datetime

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// DATE_DAYOFYEAR returns the day of year number of date.
// The return value range from 1 to 365 (366 in a leap year).
// @param {DateTime} date - Source DateTime.
// @return {Int} - A day of year number.
func DateDayOfYear(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	dt, err := runtime.CastArg[runtime.DateTime](arg, 0)

	if err != nil {
		return runtime.None, err
	}

	dayOfYear := dt.YearDay()

	return runtime.Int(dayOfYear), nil
}
