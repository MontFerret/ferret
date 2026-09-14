package math

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// max returns the largest number among the numeric elements in a list.
// @param array {Any[]} A list whose Int and Float elements are used; other elements are ignored. It is not mutated.
// @return {Float | None} The largest number as Float, or None when no numbers remain.
// @throws {TypeError} An argument has an invalid type.
func Max(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgValue(arg, 0, runtime.AssertList); err != nil {
		return runtime.None, err
	}

	return extremum(ctx, arg.(runtime.List), false)
}
