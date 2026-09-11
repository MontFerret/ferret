package arrays

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// Appends an element, optionally skipping it when it already exists.
// Existing duplicates are preserved even when unique is true.
// @param array {Any[]} Source array.
// @param value {Any} Value to append.
// @param unique {Boolean} Skip an already present value.
// @return {Any[]} New array.
// @deprecated Use arrays::append; compose arrays::unique to deduplicate the entire result.
func legacyAppend3(ctx context.Context, array, value, mode runtime.Value) (runtime.Value, error) {
	list, err := runtime.CastArg[runtime.List](array, 0)
	if err != nil {
		return runtime.None, err
	}

	unique, err := runtime.CastArg[runtime.Boolean](mode, 2)
	if err != nil {
		return runtime.None, err
	}

	next, err := copyForAppend(list)
	if err != nil {
		return runtime.None, err
	}

	if unique {
		index, err := list.IndexOf(ctx, value)
		if err != nil {
			return runtime.None, err
		}

		if index >= 0 {
			return next, nil
		}
	}

	if err := next.Append(ctx, value); err != nil {
		return runtime.None, err
	}

	return next, nil
}
