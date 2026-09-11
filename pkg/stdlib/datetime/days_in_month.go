package datetime

import (
	"context"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// DaysInMonth returns the number of days in the input's local calendar month.
// @param date {DateTime} Source date.
// @return {Int} Number of calendar days in the month, including leap days.
func DaysInMonth(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	dt, err := runtime.CastArg[runtime.DateTime](arg, 0)
	if err != nil {
		return runtime.None, err
	}

	return runtime.Int(time.Date(dt.Year(), dt.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()), nil
}
