package datetime

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// Day returns the day of date as a number.
// @param date {DateTime} Source DateTime.
// @return {Int} A day number.
func Day(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	dt, err := runtime.CastArg[runtime.DateTime](arg, 0)
	if err != nil {
		return runtime.None, err
	}

	day := dt.Day()

	return runtime.Int(day), nil
}
