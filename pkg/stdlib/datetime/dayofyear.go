package datetime

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// DATE_DAYOFYEAR returns the day of year number of date.
// The return value range from 1 to 365 (366 in a leap year).
// @param {DateTime} date - Source DateTime.
// @return {Int} - A day of year number.
func DateDayOfYear(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)
	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], types.DateTime)
	if err != nil {
		return values.None, err
	}

	dayOfYear := args[0].(values.DateTime).YearDay()

	return values.NewInt(dayOfYear), nil
}
