package datetime

import (
	"context"
	"time"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// DATE_QUARTER returns which quarter date belongs to.
// @param {DateTime} date - Source DateTime.
// @return {Int} - A quarter number.
func DateQuarter(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, 1); err != nil {
		return runtime.None, err
	}

	if err := runtime.AssertDateTime(args[0]); err != nil {
		return runtime.None, err
	}

	month := args[0].(runtime.DateTime).Month()
	quarter := runtime.NewInt(1)

	switch month {
	case time.April, time.May, time.June:
		quarter = runtime.NewInt(2)
	case time.July, time.August, time.September:
		quarter = runtime.NewInt(3)
	case time.October, time.November, time.December:
		quarter = runtime.NewInt(4)
	}

	return quarter, nil
}
