package math

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// sum returns the sum of a numeric list.
// @param numbers {Int[] | Float[]} A list containing only Int and Float values; it is not mutated.
// @return {Int | Float} The sum as Float, or integer zero for an empty list.
// @throws {TypeError} An argument or list element has an invalid type.
func Sum(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgValue(arg, 0, runtime.AssertList); err != nil {
		return runtime.None, err
	}

	sum, count, err := sumNumbers(ctx, arg.(runtime.List))
	if err != nil {
		return runtime.None, err
	}

	if count == 0 {
		return runtime.ZeroInt, nil
	}

	return runtime.Float(sum), nil
}
