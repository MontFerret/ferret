package datetime

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// DATE_SECOND returns the second of date as a number.
// @param {DateTime} date - Source DateTime.
// @return {Int} - A second number.
func DateSecond(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, 1); err != nil {
		return runtime.None, err
	}

	if err := runtime.AssertDateTime(args[0]); err != nil {
		return runtime.None, err
	}

	sec := args[0].(runtime.DateTime).Second()

	return runtime.NewInt(sec), nil
}
