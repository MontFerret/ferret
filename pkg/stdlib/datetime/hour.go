package datetime

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// DATE_HOUR returns the hour of date as a number.
// @param {DateTime} date - Source DateTime.
// @return {Int} - An hour number.
func DateHour(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, 1); err != nil {
		return runtime.None, err
	}

	if err := runtime.AssertDateTime(args[0]); err != nil {
		return runtime.None, err
	}

	hour := args[0].(runtime.DateTime).Hour()

	return runtime.NewInt(hour), nil
}
