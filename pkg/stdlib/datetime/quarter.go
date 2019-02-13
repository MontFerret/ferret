package datetime

import (
	"context"
	"time"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// DateQuarter returns which quarter date belongs to.
// @params date (DateTime) - source DateTime.
// @return (Int) - a quarter number.
func DateQuarter(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)
	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], types.DateTime)
	if err != nil {
		return values.None, err
	}

	month := args[0].(values.DateTime).Month()
	quarter := values.NewInt(1)

	switch month {
	case time.April, time.May, time.June:
		quarter = values.NewInt(2)
	case time.July, time.August, time.September:
		quarter = values.NewInt(3)
	case time.October, time.November, time.December:
		quarter = values.NewInt(4)
	}

	return quarter, nil
}
