package datetime

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// DATE_YEAR returns the year extracted from the given date.
// @param {DateTime} date - Source DateTime.
// @return {Int} - A year number.
func DateYear(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, 1); err != nil {
		return runtime.None, err
	}

	if err := runtime.AssertDateTime(args[0]); err != nil {
		return runtime.None, err
	}

	year := args[0].(runtime.DateTime).Year()

	return runtime.NewInt(year), nil
}
