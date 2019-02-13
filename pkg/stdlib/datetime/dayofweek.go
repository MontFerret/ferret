package datetime

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// DateDayOfWeek returns number of the weekday from the date. Sunday is the 0th day of week.
// @params date (DateTime) - source DateTime.
// @return (Int) - return number of the weekday.
func DateDayOfWeek(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)
	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], types.DateTime)
	if err != nil {
		return values.None, err
	}

	wday := args[0].(values.DateTime).Weekday()

	return values.NewInt(int(wday)), nil
}
