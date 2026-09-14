package arrays

import "github.com/MontFerret/ferret/v2/pkg/runtime"

// These globals are frozen migration aliases, scheduled for removal in a later
// v2 minor release. New API development belongs in arrays::. Identical contracts
// share canonical functions; mode switches and legacy-only operations stay private.
func registerLegacy(ns runtime.Namespace) {
	ns.Function().A1().
		Add("first", First).
		Add("flatten", flatten1).
		Add("last", Last).
		Add("pop", legacyPop).
		Add("shift", legacyShift).
		Add("sorted", Sorted).
		Add("sorted_unique", legacySortedUnique).
		Add("unique", Unique)
	ns.Function().A2().
		Add("append", Append).
		Add("flatten", flatten2).
		Add("nth", At).
		Add("position", legacyPosition2).
		Add("push", Append).
		Add("remove_value", Remove).
		Add("remove_nth", legacyRemoveNth).
		Add("remove_values", RemoveAny).
		Add("slice", slice2).
		Add("unshift", legacyUnshift2)
	ns.Function().A3().
		Add("append", legacyAppend3).
		Add("position", legacyPosition3).
		Add("push", legacyAppend3).
		Add("remove_value", legacyRemove3).
		Add("slice", slice3).
		Add("unshift", legacyUnshift3)
	ns.Function().Var().
		Add("intersection", Intersection).
		Add("minus", Difference).
		Add("outersection", legacyOutersection).
		Add("union", Concat).
		Add("union_distinct", Union)
}
