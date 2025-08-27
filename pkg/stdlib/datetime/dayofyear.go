package datetime

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// DATE_DAYOFYEAR returns the day of year number of date.
// The return value range from 1 to 365 (366 in a leap year).
// @param {DateTime} date - Source DateTime.
// @return {Int} - A day of year number.
func DateDayOfYear(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, 1); err != nil {
		return runtime.None, err
	}

	if err := runtime.AssertDateTime(args[0]); err != nil {
		return runtime.None, err
	}

	dayOfYear := args[0].(runtime.DateTime).YearDay()

	return runtime.NewInt(dayOfYear), nil
}
