package math

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// median returns the middle numeric value, averaging the two middle numbers for even numeric counts.
// @param array {Any[]} A list whose Int and Float elements are used; other elements are ignored. It is not mutated.
// @return {Int | Float | None} The selected middle value or their Float mean, or None when no numbers remain.
// @throws {TypeError} An argument has an invalid type.
func Median(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgValue(arg, 0, runtime.AssertList); err != nil {
		return runtime.None, err
	}

	values, err := snapshotNumbers(ctx, arg.(runtime.List))
	if err != nil {
		return runtime.None, err
	}

	if len(values) == 0 {
		return runtime.None, nil
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
