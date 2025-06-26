package datetime

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// DATE_MONTH returns the month of date as a number.
// @param {DateTime} date - Source DateTime.
// @return {Int} - A month number.
func DateMonth(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, 1); err != nil {
		return runtime.None, err
	}

	if err := runtime.AssertDateTime(args[0]); err != nil {
		return runtime.None, err
	}

	month := args[0].(runtime.DateTime).Month()

	return runtime.NewInt(int(month)), nil
}
