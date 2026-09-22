package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// ceil returns the least integer value greater than or equal to a given value.
// @param number {Int | Float} Input number.
// @return {Int} The least integer value greater than or equal to a given value.
// @throws {RangeError} The rounded result is not a finite signed 64-bit integer.
func Ceil(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgValue(arg, 0, runtime.AssertNumber); err != nil {
		return runtime.None, err
	}

	if integer, ok := arg.(runtime.Int); ok {
		return integer, nil
	}

	return roundedInt(math.Ceil(float64(arg.(runtime.Float))))
}
