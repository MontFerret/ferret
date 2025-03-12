package datetime

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// DATE_MONTH returns the month of date as a number.
// @param {DateTime} date - Source DateTime.
// @return {Int} - A month number.
func DateMonth(_ context.Context, args ...core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 1, 1); err != nil {
		return core.None, err
	}

	if err := core.AssertDateTime(args[0]); err != nil {
		return core.None, err
	}

	month := args[0].(core.DateTime).Month()

	return core.NewInt(int(month)), nil
}
