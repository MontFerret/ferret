package objects

import "github.com/MontFerret/ferret/pkg/runtime/core"

func RegisterLib(ns core.Namespace) error {
	return ns.RegisterFunctions(core.Functions{
		"HAS":             Has,
		"KEYS":            Keys,
		"KEEP_KEYS":       KeepKeys,
		"MERGE":           Merge,
		"ZIP":             Zip,
		"VALUES":          Values,
		"MERGE_RECURSIVE": MergeRecursive,
	})
}
