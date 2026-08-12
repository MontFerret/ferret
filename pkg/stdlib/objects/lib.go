package objects

import "github.com/MontFerret/ferret/v2/pkg/runtime"

func RegisterLib(ns runtime.Namespace) {
	ns.Function().A1().
		Add("keys", keys1).
		Add("values", Values)
	ns.Function().A2().
		Add("has", Has).
		Add("keys", keys2).
		Add("zip", Zip)
	ns.Function().Var().
		Add("keep_keys", KeepKeys).
		Add("merge", Merge).
		Add("merge_recursive", MergeRecursive)
}
