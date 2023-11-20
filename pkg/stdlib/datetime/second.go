package datetime

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// DATE_SECOND returns the second of date as a number.
// @param {DateTime} date - Source DateTime.
// @return {Int} - A second number.
func DateSecond(_ context.Context, args ...core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 1, 1); err != nil {
		return values.None, err
	}

	if err := values.AssertDateTime(args[0]); err != nil {
		return values.None, err
	}

	sec := args[0].(values.DateTime).Second()

	return values.NewInt(sec), nil
}
