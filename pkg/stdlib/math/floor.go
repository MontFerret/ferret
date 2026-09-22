package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// floor returns the greatest integer value less than or equal to a given value.
// @param number {Int | Float} Input number.
// @return {Int} The greatest integer value less than or equal to a given value.
// @throws {RangeError} The rounded result is not a finite signed 64-bit integer.
func Floor(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgValue(arg, 0, runtime.AssertNumber); err != nil {
		return runtime.None, err
	}

	if integer, ok := arg.(runtime.Int); ok {
		return integer, nil
	}

	return roundedInt(math.Floor(float64(arg.(runtime.Float))))
}
