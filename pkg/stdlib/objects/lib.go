package objects

import "github.com/MontFerret/ferret/pkg/runtime"

func RegisterLib(ns runtime.Namespace) error {
	ns.Functions().
		Set("HAS", Has).
		Set("KEYS", Keys).
		Set("KEEP_KEYS", KeepKeys).
		Set("MERGE", Merge).
		Set("ZIP", Zip).
		Set("VALUES", Values).
		Set("MERGE_RECURSIVE", MergeRecursive)

	return nil
}
