package objects

import "github.com/MontFerret/ferret/v2/pkg/runtime"

func RegisterLib(ns runtime.Namespace) {
	ns.Function().A1().
		Add("KEYS", keys1).
		Add("VALUES", Values)
	ns.Function().A2().
		Add("HAS", Has).
		Add("KEYS", keys2).
		Add("ZIP", Zip)
	ns.Function().Var().
		Add("KEEP_KEYS", KeepKeys).
		Add("MERGE", Merge).
		Add("MERGE_RECURSIVE", MergeRecursive)
}
