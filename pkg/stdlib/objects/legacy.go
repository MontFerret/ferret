package objects

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// These temporary globals share canonical semantics. Separate forwarding
// declarations attach deprecation metadata without deprecating object:: functions.
func registerLegacy(ns runtime.Namespace) {
	ns.Function().A1().
		Add("keys", legacyKeys).
		Add("values", legacyValues)
	ns.Function().A2().
		Add("has", legacyHas).
		Add("zip", legacyZip)
	ns.Function().Var().
		Add("keep_keys", legacyKeepKeys).
		Add("merge", legacyMerge).
		Add("merge_recursive", legacyMergeRecursive)
}

// keys returns an independent list of map keys in unspecified order.
// @param value {Map} Map whose keys are returned.
// @return {String[]} Independent list of keys in unspecified order.
// @deprecated Use object::keys instead.
func legacyKeys(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	return Keys(ctx, arg)
}

// values returns cloned or copied map values in unspecified order.
// @param value {Map} Map whose values are returned.
// @return {Any[]} Independent list of values, following each value's clone or copy contract.
// @deprecated Use object::values instead.
func legacyValues(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	return Values(ctx, arg)
}

// has checks key existence, including keys whose value is none.
// @param object {Map} Object or map to inspect.
// @param key {String} The key name string.
// @return {Boolean} True if the key exists else false.
// @deprecated Use object::has_key instead.
func legacyHas(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	return HasKey(ctx, arg1, arg2)
}

// keep_keys clones the source and retains only the requested keys.
// Missing and repeated keys are ignored.
// @param value {Map} Source map.
// @param keys {String|String[], repeated} Variadic keys, or one list of keys; an empty list keeps nothing.
// @return {Map} Independent map containing only selected keys.
// @deprecated Use object::keep_keys instead.
func legacyKeepKeys(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	return KeepKeys(ctx, args...)
}

// merge copies maps into an independent destination; later values replace earlier values.
// @param values {Map|Map[], repeated} Variadic maps, or one list of maps.
// @return {Map} Shallow merge with cloned or copied values.
// @deprecated Use object::merge instead.
func legacyMerge(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	return Merge(ctx, args...)
}

// merge_recursive copies maps into an independent destination, merging nested maps recursively.
// Arrays, scalars, and none are replaced by later values.
// @param values {Map|Map[], repeated} Variadic maps, or one list of maps.
// @return {Map} Deep merge with cloned or copied values.
// @deprecated Use object::merge_deep instead.
func legacyMergeRecursive(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	return MergeDeep(ctx, args...)
}

// zip constructs an object from parallel key and value lists of equal length.
// Keys must be strings. Duplicate keys use the last value in the input.
// @param keys {String[]} List of object keys.
// @param values {Any[]} Parallel list of values.
// @return {Map} Independent object containing cloned or copied values.
// @deprecated Use object::zip instead.
func legacyZip(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	return Zip(ctx, arg1, arg2)
}
