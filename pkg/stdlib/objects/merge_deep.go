package objects

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// MergeDeep copies maps into an independent destination, merging nested maps recursively.
// Arrays, scalars, and none are replaced by later values.
// @param values {Map|Map[], repeated} Variadic maps, or one list of maps.
// @return {Map} Deep merge with cloned or copied values.
func MergeDeep(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	return mergeObjects(ctx, args, runtime.MergeMapsDeepInto)
}
