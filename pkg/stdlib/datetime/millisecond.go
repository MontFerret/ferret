package datetime

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// DATE_MILLISECOND returns the millisecond of date as a number.
// @param {DateTime} date - Source DateTime.
// @return {Int} - A millisecond number.
func DateMillisecond(_ context.Context, args ...core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 1, 1); err != nil {
		return values.None, err
	}

	if err := values.AssertDateTime(args[0]); err != nil {
		return values.None, err
	}

	msec := args[0].(values.DateTime).Nanosecond() / 1000000

	return values.NewInt(msec), nil
}
