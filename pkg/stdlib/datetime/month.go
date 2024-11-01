package datetime

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// DATE_MONTH returns the month of date as a number.
// @param {DateTime} date - Source DateTime.
// @return {Int} - A month number.
func DateMonth(_ context.Context, args ...core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 1, 1); err != nil {
		return values.None, err
	}

	if err := values.AssertDateTime(args[0]); err != nil {
		return values.None, err
	}

	month := args[0].(values.DateTime).Month()

	return values.NewInt(int(month)), nil
}
