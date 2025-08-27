package datetime

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// DATE_DAY returns the day of date as a number.
// @param {DateTime} date - Source DateTime.
// @return {Int} - A day number.
func DateDay(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, 1); err != nil {
		return runtime.None, err
	}

	if err := runtime.AssertDateTime(args[0]); err != nil {
		return runtime.None, err
	}

	day := args[0].(runtime.DateTime).Day()

	return runtime.NewInt(day), nil
}
