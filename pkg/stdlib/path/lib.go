package path

import "github.com/MontFerret/ferret/pkg/runtime/core"

// RegisterLib register `PATH` namespace functions.
// @namespace PATH
func RegisterLib(ns core.Namespace) error {
	return ns.
		Namespace("PATH").
		RegisterFunctions(
			core.NewFunctionsFromMap(map[string]core.Function{
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
