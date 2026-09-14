package arrays

import "github.com/MontFerret/ferret/v2/pkg/runtime"

// RegisterLib registers arrays, explicit arrays::mut operations, and legacy globals.
// @namespace arrays
func RegisterLib(ns runtime.Namespace) {
	canonical := ns.Namespace("arrays")
	canonical.Function().A1().
		Add("first", First).
		Add("last", Last).
		Add("flatten", flatten1).
		Add("unique", Unique).
		Add("sorted", Sorted)
	canonical.Function().A2().
		Add("range", range2).
		Add("at", At).
		Add("append", Append).
		Add("flatten", flatten2).
		Add("slice", slice2).
		Add("contains", Contains).
		Add("index_of", IndexOf).
		Add("remove", Remove).
		Add("remove_at", RemoveAt).
		Add("remove_any", RemoveAny)
	canonical.Function().A3().
		Add("range", range3).
		Add("slice", slice3)
	canonical.Function().Var().
		Add("concat", Concat).
		Add("union", Union).
		Add("intersection", Intersection).
		Add("difference", Difference).
		Add("symmetric_difference", SymmetricDifference)

	mutable := canonical.Namespace("mut")
	mutable.Function().A1().
		Add("pop", PopMutable).
		Add("shift", ShiftMutable).
		Add("clear", ClearMutable).
		Add("sort", SortMutable)
	mutable.Function().A2().
		Add("push", PushMutable).
		Add("unshift", UnshiftMutable).
		Add("remove", RemoveMutable).
		Add("remove_at", RemoveAtMutable)
	mutable.Function().A3().
		Add("set", SetMutable).
		Add("insert", InsertMutable)

	registerLegacy(ns)
}
