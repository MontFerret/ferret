package math

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// min returns the smallest number among the numeric elements in a list.
// @param array {Any[]} A list whose Int and Float elements are used; other elements are ignored. It is not mutated.
// @return {Float | None} The smallest number as Float, or None when no numbers remain.
// @throws {TypeError} An argument has an invalid type.
// @deprecated Use math::min instead.
func Min(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgValue(arg, 0, runtime.AssertList); err != nil {
		return runtime.None, err
	}

	return extremum(ctx, arg.(runtime.List), true, filterNonNumbers)
}

// min returns the smallest number in a numeric list.
// @param values {Int[] | Float[]} Numeric list; non-numeric elements are rejected. The list is not mutated.
// @return {Float | None} The minimum as Float, or None for an empty list.
// @throws {TypeError} An argument or list element has an invalid type.
func canonicalMin(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgValue(arg, 0, runtime.AssertList); err != nil {
		return runtime.None, err
	}

	return extremum(ctx, arg.(runtime.List), true, strictNumbers)
}
