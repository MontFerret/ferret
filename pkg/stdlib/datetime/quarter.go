package datetime

import (
	"context"
	"time"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// DATE_QUARTER returns which quarter date belongs to.
// @param {DateTime} date - Source DateTime.
// @return {Int} - A quarter number.
func DateQuarter(_ context.Context, args ...core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 1, 1); err != nil {
		return core.None, err
	}

	if err := core.AssertDateTime(args[0]); err != nil {
		return core.None, err
	}

	month := args[0].(core.DateTime).Month()
	quarter := core.NewInt(1)

	switch month {
	case time.April, time.May, time.June:
		quarter = core.NewInt(2)
	case time.July, time.August, time.September:
		quarter = core.NewInt(3)
	case time.October, time.November, time.December:
		quarter = core.NewInt(4)
	}

	return quarter, nil
}
