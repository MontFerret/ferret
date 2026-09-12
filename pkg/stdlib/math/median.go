package math

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// median returns the middle value of a numeric list, averaging the two middle values for even lengths.
// @param array {Int[] | Float[]} A list containing only Int and Float values; it is not mutated.
// @return {Int | Float} The selected middle value or their Float mean, or NaN for an empty list.
// @throws {TypeError} An argument or list element has an invalid type.
func Median(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgValue(arg, 0, runtime.AssertList); err != nil {
		return runtime.None, err
	}

	values, err := snapshotNumbers(ctx, arg.(runtime.List))
	if err != nil {
		return runtime.None, err
	}

	if len(values) == 0 {
		return runtime.NaN(), nil
	}

	if err := sortNumbers(ctx, values, false); err != nil {
		return runtime.None, err
	}

	middle := len(values) / 2
	if len(values)%2 != 0 {
		return values[middle], nil
	}

	return runtime.Float((toFloat(values[middle-1]) + toFloat(values[middle])) / 2), nil
}
