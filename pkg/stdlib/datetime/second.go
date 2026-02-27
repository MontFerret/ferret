package datetime

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// DATE_SECOND returns the second of date as a number.
// @param {DateTime} date - Source DateTime.
// @return {Int} - A second number.
func DateSecond(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	dt, err := runtime.CastArg[runtime.DateTime](arg, 0)

	if err != nil {
		return runtime.None, err
	}

	sec := dt.Second()

	return runtime.Int(sec), nil
}
