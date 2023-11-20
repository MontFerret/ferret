package datetime

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// DATE_YEAR returns the year extracted from the given date.
// @param {DateTime} date - Source DateTime.
// @return {Int} - A year number.
func DateYear(_ context.Context, args ...core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 1, 1); err != nil {
		return values.None, err
	}

	if err := values.AssertDateTime(args[0]); err != nil {
		return values.None, err
	}

	year := args[0].(values.DateTime).Year()

	return values.NewInt(year), nil
}
