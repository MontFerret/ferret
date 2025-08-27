package objects

import "github.com/MontFerret/ferret/pkg/runtime"

func RegisterLib(ns runtime.Namespace) error {
	return ns.RegisterFunctions(
		runtime.NewFunctionsFromMap(map[string]runtime.Function{
			"HAS":             Has,
			"KEYS":            Keys,
			"KEEP_KEYS":       KeepKeys,
			"MERGE":           Merge,
			"ZIP":             Zip,
			"VALUES":          Values,
			"MERGE_RECURSIVE": MergeRecursive,
		}))
}
