package datetime

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// DATE_MILLISECOND returns the millisecond of date as a number.
// @param {DateTime} date - Source DateTime.
// @return {Int} - A millisecond number.
func DateMillisecond(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, 1); err != nil {
		return runtime.None, err
	}

	if err := runtime.AssertDateTime(args[0]); err != nil {
		return runtime.None, err
	}

	msec := args[0].(runtime.DateTime).Nanosecond() / 1000000

	return runtime.NewInt(msec), nil
}
