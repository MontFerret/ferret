package arrays

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// Removes equal values, optionally limiting the number of removals.
// @param array {Any[]} Source array.
// @param value {Any} Value to remove.
// @param limit {Int} Maximum removals; negative is unlimited and zero removes nothing.
// @return {Any[]} New array retaining the other elements in order.
// @deprecated Use arrays::remove to remove all matches. Limits remain legacy-only.
func legacyRemove3(ctx context.Context, array, value, count runtime.Value) (runtime.Value, error) {
	list, err := runtime.CastArg[runtime.List](array, 0)
	if err != nil {
		return runtime.None, err
	}

	limit, err := runtime.CastArg[runtime.Int](count, 2)
	if err != nil {
		return runtime.None, err
	}

	return removeLimited(ctx, list, value, limit)
}

// Removes an indexed value from a copy, preserving the host list's removal policy.
// @param array {Any[]} Source array.
// @param index {Int} Position to remove.
// @return {Any[]} New array without the indexed element.
// @deprecated Use arrays::remove_at for bounds-checked removal.
func legacyRemoveNth(ctx context.Context, array, index runtime.Value) (runtime.Value, error) {
	list, position, err := runtime.CastArgs2[runtime.List, runtime.Int](array, index)
	if err != nil {
		return runtime.None, err
	}

	next, err := runtime.Copy(list)
	if err != nil {
		return runtime.None, err
	}

	_, err = next.RemoveAt(ctx, position)

	return next, err
}
