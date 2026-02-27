package datetime

import (
	"context"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// DATE_QUARTER returns which quarter date belongs to.
// @param {DateTime} date - Source DateTime.
// @return {Int} - A quarter number.
func DateQuarter(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	dt, err := runtime.CastArg[runtime.DateTime](arg, 0)

	if err != nil {
		return runtime.None, err
	}

	month := dt.Month()
	quarter := runtime.Int(1)

	switch month {
	case time.April, time.May, time.June:
		quarter = runtime.NewInt(2)
	case time.July, time.August, time.September:
		quarter = runtime.NewInt(3)
	case time.October, time.November, time.December:
		quarter = runtime.NewInt(4)
	default:
		quarter = runtime.NewInt(1)
	}

	return quarter, nil
}
