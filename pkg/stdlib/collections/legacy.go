package collections

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// Separate forwarding declarations attach deprecation metadata only to globals.
func registerLegacy(ns runtime.Namespace) {
	ns.Function().A1().
		Add("count", legacyCount).
		Add("count_distinct", legacyCountDistinct).
		Add("reverse", legacyReverse)

	ns.Function().A2().
		Add("includes", legacyIncludes)
}

// count returns the measured length or counts yielded values in one traversal.
// @param collection {Iterable} Source whose yielded values are counted.
// @return {Int} Number of yielded values, or the source's measured length.
// @deprecated Use collections::count instead.
func legacyCount(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	return Count(ctx, arg)
}

// count_distinct counts distinct yielded values using canonical runtime equality.
// @param collection {Iterable} Source whose distinct yielded values are counted.
// @return {Int} Number of distinct yielded values.
// @deprecated Use collections::count_distinct instead.
func legacyCountDistinct(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	return CountDistinct(ctx, arg)
}

// includes checks membership using host capabilities or an iterable scan.
// @param haystack {String | Containable | Iterable} The value container.
// @param needle {Any} The target value to assert.
// @return {Boolean} A boolean value that indicates whether a container contains a given value.
// @deprecated Use collections::includes instead.
func legacyIncludes(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	return Includes(ctx, arg1, arg2)
}

// reverse reverses Unicode code points or indexed list elements, preserving the list family.
// @param value {String | List} The string or indexed list to reverse.
// @return {String | List} A reversed string or new list of the source's implementation family.
// @deprecated Use collections::reverse instead.
func legacyReverse(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	return Reverse(ctx, arg)
}
