package datetime

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// DATE_LEAPYEAR returns true if date is in a leap year else false.
// @param {DateTime} date - Source DateTime.
// @return {Boolean} - Date is in a leap year.
func DateLeapYear(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, 1); err != nil {
		return runtime.None, err
	}

	if err := runtime.AssertDateTime(args[0]); err != nil {
		return runtime.None, err
	}

	year := args[0].(runtime.DateTime).Year()

	return runtime.NewBoolean(isLeap(year)), nil
}

func isLeap(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}
