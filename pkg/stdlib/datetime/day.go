package datetime

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// DATE_DAY returns the day of date as a number.
// @param {DateTime} date - Source DateTime.
// @return {Int} - A day number.
func DateDay(_ context.Context, args ...core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 1, 1); err != nil {
		return core.None, err
	}

	if err := core.AssertDateTime(args[0]); err != nil {
		return core.None, err
	}

	day := args[0].(core.DateTime).Day()

	return core.NewInt(day), nil
}
