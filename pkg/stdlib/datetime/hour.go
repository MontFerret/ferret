package datetime

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// DATE_HOUR returns the hour of date as a number.
// @param {DateTime} date - Source DateTime.
// @return {Int} - An hour number.
func DateHour(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	dt, err := runtime.CastArg[runtime.DateTime](arg, 0)

	if err != nil {
		return runtime.None, err
	}

	hour := dt.Hour()

	return runtime.Int(hour), nil
}
