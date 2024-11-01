package datetime

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// DATE_DAY returns the day of date as a number.
// @param {DateTime} date - Source DateTime.
// @return {Int} - A day number.
func DateDay(_ context.Context, args ...core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 1, 1); err != nil {
		return values.None, err
	}

	if err := values.AssertDateTime(args[0]); err != nil {
		return values.None, err
	}

	day := args[0].(values.DateTime).Day()

	return values.NewInt(day), nil
}
