package math

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// min returns the smallest number in a numeric list.
// @param array {Int[] | Float[]} A list containing only Int and Float values; it is not mutated.
// @return {Float | None} The smallest number as Float, or None for an empty list.
// @throws {TypeError} An argument or list element has an invalid type.
func Min(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgValue(arg, 0, runtime.AssertList); err != nil {
		return runtime.None, err
	}

	return extremum(ctx, arg.(runtime.List), true)
}
