package objects

import "github.com/MontFerret/ferret/v2/pkg/runtime"

// RegisterLib registers immutable object functions in the object namespace.
// @namespace object
func RegisterLib(ns runtime.Namespace) {
	ns = ns.Namespace("object")

	ns.Function().A1().
		Add("keys", Keys).
		Add("values", Values).
		Add("entries", Entries).
		Add("from_entries", FromEntries)
	ns.Function().A2().
		Add("has_key", HasKey).
		Add("zip", Zip)
	ns.Function().Var().
		Add("keep_keys", KeepKeys).
		Add("omit_keys", OmitKeys).
		Add("merge", Merge).
		Add("merge_deep", MergeDeep)
}
