package datetime

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// DATE_SECOND returns the second of date as a number.
// @param {DateTime} date - Source DateTime.
// @return {Int} - A second number.
func DateSecond(_ context.Context, args ...core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 1, 1); err != nil {
		return core.None, err
	}

	if err := core.AssertDateTime(args[0]); err != nil {
		return core.None, err
	}

	sec := args[0].(core.DateTime).Second()

	return core.NewInt(sec), nil
}
