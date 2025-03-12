package datetime

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// DATE_HOUR returns the hour of date as a number.
// @param {DateTime} date - Source DateTime.
// @return {Int} - An hour number.
func DateHour(_ context.Context, args ...core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 1, 1); err != nil {
		return core.None, err
	}

	if err := core.AssertDateTime(args[0]); err != nil {
		return core.None, err
	}

	hour := args[0].(core.DateTime).Hour()

	return core.NewInt(hour), nil
}
