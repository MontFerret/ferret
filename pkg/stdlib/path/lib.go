package path

import "github.com/MontFerret/ferret/pkg/runtime"

// RegisterLib register `PATH` namespace functions.
// @namespace PATH
func RegisterLib(ns runtime.Namespace) error {
	return ns.
		Namespace("PATH").
		RegisterFunctions(
			runtime.NewFunctionsFromMap(map[string]runtime.Function{
				"BASE":     Base,
				"CLEAN":    Clean,
				"DIR":      Dir,
				"EXT":      Ext,
				"IS_ABS":   IsAbs,
				"JOIN":     Join,
				"MATCH":    Match,
				"SEPARATE": Separate,
			}))
}
