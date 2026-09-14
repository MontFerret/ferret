package math

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// max returns the largest number among the numeric elements in a list.
// @param array {Any[]} A list whose Int and Float elements are used; other elements are ignored. It is not mutated.
// @return {Float | None} The largest number as Float, or None when no numbers remain.
// @throws {TypeError} An argument has an invalid type.
// @deprecated Use math::max instead.
func Max(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgValue(arg, 0, runtime.AssertList); err != nil {
		return runtime.None, err
	}

	return extremum(ctx, arg.(runtime.List), false, filterNonNumbers)
}

// max returns the largest number in a numeric list.
// @param values {Int[] | Float[]} Numeric list; non-numeric elements are rejected. The list is not mutated.
// @return {Float | None} The maximum as Float, or None for an empty list.
// @throws {TypeError} An argument or list element has an invalid type.
func canonicalMax(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgValue(arg, 0, runtime.AssertList); err != nil {
		return runtime.None, err
	}

	return extremum(ctx, arg.(runtime.List), false, strictNumbers)
}
