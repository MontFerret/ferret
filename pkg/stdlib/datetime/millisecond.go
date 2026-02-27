package datetime

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// DATE_MILLISECOND returns the millisecond of date as a number.
// @param {DateTime} date - Source DateTime.
// @return {Int} - A millisecond number.
func DateMillisecond(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	dt, err := runtime.CastArg[runtime.DateTime](arg, 0)

	if err != nil {
		return runtime.None, err
	}

	msec := dt.Nanosecond() / 1000000

	return runtime.Int(msec), nil
}
