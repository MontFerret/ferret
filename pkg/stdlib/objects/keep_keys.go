package objects

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// KeepKeys clones the source and retains only the requested keys.
// Missing and repeated keys are ignored.
// @param value {Map} Source map.
// @param keys {String|String[], repeated} Variadic keys, or one list of keys; an empty list keeps nothing.
// @return {Map} Independent map containing only selected keys.
func KeepKeys(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	return filterKeys(ctx, args, runtime.KeepMapKeys)
}

// KeepKeysMutable removes unselected keys from the original target and returns it.
// Retained values are untouched. Missing and repeated keys are ignored; errors
// may leave earlier removals applied.
// @param target {Map} Map to mutate.
// @param keys {String|String[], repeated} Variadic keys, or one list of keys; an empty list keeps nothing.
// @return {Map} The original target containing only selected keys.
func KeepKeysMutable(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	return filterMutableKeys(ctx, args, runtime.KeepMapKeys)
}
