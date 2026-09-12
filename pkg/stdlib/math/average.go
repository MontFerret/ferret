package math

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// average returns the arithmetic mean of a numeric list.
// @param array {Int[] | Float[]} A list containing only Int and Float values; it is not mutated.
// @return {Float} The arithmetic mean, or NaN for an empty list.
// @throws {TypeError} An argument or list element has an invalid type.
func Average(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgValue(arg, 0, runtime.AssertList); err != nil {
		return runtime.None, err
	}

	result, _, err := mean(ctx, arg.(runtime.List))
	if err != nil {
		return runtime.None, err
	}

	return result, nil
}
