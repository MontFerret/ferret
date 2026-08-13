package datetime

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// date_month returns the month of date as a number.
// @param date {DateTime} Source DateTime.
// @return {Int} A month number.
func DateMonth(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	dt, err := runtime.CastArg[runtime.DateTime](arg, 0)

	if err != nil {
		return runtime.None, err
	}

	month := dt.Month()

	return runtime.Int(month), nil
}
