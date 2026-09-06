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
