package datetime

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// DATE_YEAR returns the year extracted from the given date.
// @param {DateTime} date - Source DateTime.
// @return {Int} - A year number.
func DateYear(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	dt, err := runtime.CastArg[runtime.DateTime](arg, 0)

	if err != nil {
		return runtime.None, err
	}

	year := dt.Year()

	return runtime.Int(year), nil
}
