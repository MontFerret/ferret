package datetime

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// IsLeapYear returns true if date is in a leap year else false.
// @param date {DateTime} Source DateTime.
// @return {Boolean} Whether the date is in a leap year.
func IsLeapYear(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	dt, err := runtime.CastArg[runtime.DateTime](arg, 0)
	if err != nil {
		return runtime.None, err
	}

	year := dt.Year()

	return runtime.Boolean(isLeap(year)), nil
}

func isLeap(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}
