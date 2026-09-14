package math

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// median returns the middle numeric value, averaging the two middle numbers for even numeric counts.
// @param array {Any[]} A list whose Int and Float elements are used; other elements are ignored. It is not mutated.
// @return {Int | Float | None} The selected middle value or their Float mean, or None when no numbers remain.
// @throws {TypeError} An argument has an invalid type.
// @deprecated Use math::median instead.
func Median(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	return median(ctx, arg, filterNonNumbers)
}

// median returns the central number, averaging the middle pair for even counts.
// @param values {Int[] | Float[]} Numeric list; non-numeric elements are rejected. The list is not mutated.
// @return {Int | Float} The selected number or Float mean, or NaN for an empty list.
// @throws {TypeError} An argument or list element has an invalid type.
func canonicalMedian(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	value, err := median(ctx, arg, strictNumbers)
	if err != nil {
		return runtime.None, err
	}

	if value == runtime.None {
		return runtime.NaN(), nil
	}

	return value, nil
}

func median(ctx context.Context, arg runtime.Value, policy numberPolicy) (runtime.Value, error) {
	if err := runtime.ValidateArgValue(arg, 0, runtime.AssertList); err != nil {
		return runtime.None, err
	}

	values, err := snapshotNumbers(ctx, arg.(runtime.List), policy)
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
