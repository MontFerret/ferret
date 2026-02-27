package datetime

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// DATE_DAY returns the day of date as a number.
// @param {DateTime} date - Source DateTime.
// @return {Int} - A day number.
func DateDay(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	dt, err := runtime.CastArg[runtime.DateTime](arg, 0)

	if err != nil {
		return runtime.None, err
	}

	day := dt.Day()

	return runtime.Int(day), nil
}
