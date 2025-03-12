package datetime

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// DATE_YEAR returns the year extracted from the given date.
// @param {DateTime} date - Source DateTime.
// @return {Int} - A year number.
func DateYear(_ context.Context, args ...core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 1, 1); err != nil {
		return core.None, err
	}

	if err := core.AssertDateTime(args[0]); err != nil {
		return core.None, err
	}

	year := args[0].(core.DateTime).Year()

	return core.NewInt(year), nil
}
