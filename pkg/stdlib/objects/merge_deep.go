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

// MergeDeepMutable deep merges sources into the original target, returning it.
// Conflicting nested maps are cloned before modification so external nested
// aliases remain unchanged. Arrays, scalars, and none are replaced. Errors may
// leave earlier updates applied. A target without sources is returned unchanged.
// @param target {Map} Map to mutate; its top-level aliases observe changes.
// @param sources {Map|Map[], repeated} Zero or more maps, or one list of maps.
// @return {Map} The original target with independently merged nested branches.
func MergeDeepMutable(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	return mergeMutableObjects(ctx, args, runtime.MergeMapsDeepIsolatedInto)
}
