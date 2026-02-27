package datetime

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// DATE_MINUTE returns the minute of date as a number.
// @param {DateTime} date -Source DateTime.
// @return {Int} - A minute number.
func DateMinute(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	dt, err := runtime.CastArg[runtime.DateTime](arg, 0)

	if err != nil {
		return runtime.None, err
	}

	mmin := dt.Minute()

	return runtime.Int(mmin), nil
}
