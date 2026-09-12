package math

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// max returns the largest number in a numeric list.
// @param array {Int[] | Float[]} A list containing only Int and Float values; it is not mutated.
// @return {Float | None} The largest number as Float, or None for an empty list.
// @throws {TypeError} An argument or list element has an invalid type.
func Max(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgValue(arg, 0, runtime.AssertList); err != nil {
		return runtime.None, err
	}

	return extremum(ctx, arg.(runtime.List), false)
}
