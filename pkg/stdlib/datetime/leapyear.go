package datetime

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// DATE_LEAPYEAR returns true if date is in a leap year else false.
// @param {DateTime} date - Source DateTime.
// @return {Boolean} - Date is in a leap year.
func DateLeapYear(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)
	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], types.DateTime)
	if err != nil {
		return values.None, err
	}

	year := args[0].(values.DateTime).Year()

	return values.NewBoolean(isLeap(year)), nil
}

func isLeap(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}
