package math

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// sum returns the sum of the numeric elements in a list.
// @param numbers {Any[]} A list whose Int and Float elements are used; other elements are ignored. It is not mutated.
// @return {Int | Float} The sum as Float, including zero when all elements are ignored; integer zero for an empty list.
// @throws {TypeError} An argument has an invalid type.
// @deprecated Use math::sum instead.
func Sum(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	return sum(ctx, arg, filterNonNumbers)
}

// sum returns the sum of a numeric list.
// @param values {Int[] | Float[]} Numeric list; non-numeric elements are rejected. The list is not mutated.
// @return {Int | Float} The sum as Float, or integer zero for an empty list.
// @throws {TypeError} An argument or list element has an invalid type.
func canonicalSum(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	return sum(ctx, arg, strictNumbers)
}

func sum(ctx context.Context, arg runtime.Value, policy numberPolicy) (runtime.Value, error) {
	if err := runtime.ValidateArgValue(arg, 0, runtime.AssertList); err != nil {
		return runtime.None, err
	}

	sum, counts, err := sumNumbers(ctx, arg.(runtime.List), policy)
	if err != nil {
		return runtime.None, err
	}

	if counts.total == 0 {
		return runtime.ZeroInt, nil
	}

	return runtime.Float(sum), nil
}
