package path

import "github.com/MontFerret/ferret/v2/pkg/runtime"

// RegisterLib register `PATH` namespace functions.
// @namespace PATH
func RegisterLib(ns runtime.Namespace) {
	ns.Function().A1().
		Add("BASE", Base).
		Add("CLEAN", Clean).
		Add("DIR", Dir).
		Add("EXT", Ext).
		Add("IS_ABS", IsAbs).
		Add("SEPARATE", Separate)

	ns.Function().A2().
		Add("MATCH", Match)

	ns.Function().Var().
		Add("JOIN", Join)
}
