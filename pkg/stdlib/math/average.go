package math

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// average returns the arithmetic mean of the numeric elements in a list.
// @param array {Any[]} A list whose Int and Float elements are used; other elements are ignored. It is not mutated.
// @return {Float} The arithmetic mean, or zero when no numbers remain.
// @throws {TypeError} An argument has an invalid type.
// @deprecated Use math::mean instead.
func Average(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgValue(arg, 0, runtime.AssertList); err != nil {
		return runtime.None, err
	}

	result, _, err := mean(ctx, arg.(runtime.List), filterNonNumbers)
	if err != nil {
		return runtime.None, err
	}

	return result, nil
}
