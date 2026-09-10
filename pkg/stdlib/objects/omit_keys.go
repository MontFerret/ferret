package objects

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// OmitKeys clones the source and removes the requested keys.
// Missing and repeated keys are ignored.
// @param value {Map} Source map.
// @param keys {String|String[], repeated} Variadic keys, or one list of keys; an empty list omits nothing.
// @return {Map} Independent map without selected keys.
func OmitKeys(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	return filterKeys(ctx, args, runtime.OmitMapKeys)
}

// OmitKeysMutable removes selected keys from the original target and returns it.
// Retained values are untouched. Missing and repeated keys are ignored; errors
// may leave earlier removals applied.
// @param target {Map} Map to mutate.
// @param keys {String|String[], repeated} Variadic keys, or one list of keys; an empty list omits nothing.
// @return {Map} The original target without selected keys.
func OmitKeysMutable(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	return filterMutableKeys(ctx, args, runtime.OmitMapKeys)
}
