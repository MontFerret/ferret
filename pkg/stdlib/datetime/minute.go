package datetime

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// DATE_MINUTE returns the minute of date as a number.
// @param {DateTime} date -Source DateTime.
// @return {Int} - A minute number.
func DateMinute(_ context.Context, args ...core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 1, 1); err != nil {
		return core.None, err
	}

	if err := core.AssertDateTime(args[0]); err != nil {
		return core.None, err
	}

	min := args[0].(core.DateTime).Minute()

	return core.NewInt(min), nil
}
