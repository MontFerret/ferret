package objects

import "github.com/MontFerret/ferret/v2/pkg/runtime"

// RegisterLib registers object functions, explicit object::mut operations, and deprecated globals.
// @namespace object
func RegisterLib(ns runtime.Namespace) {
	canonical := ns.Namespace("object")

	canonical.Function().A1().
		Add("keys", Keys).
		Add("values", Values).
		Add("entries", Entries).
		Add("from_entries", FromEntries)
	canonical.Function().A2().
		Add("has_key", HasKey).
		Add("zip", Zip)
	canonical.Function().Var().
		Add("keep_keys", KeepKeys).
		Add("omit_keys", OmitKeys).
		Add("merge", Merge).
		Add("merge_deep", MergeDeep)

	canonical.Namespace("mut").Function().Var().
		Add("keep_keys", KeepKeysMutable).
		Add("omit_keys", OmitKeysMutable).
		Add("merge", MergeMutable).
		Add("merge_deep", MergeDeepMutable)

	registerLegacy(ns)
}
