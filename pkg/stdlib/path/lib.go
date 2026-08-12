package path

import "github.com/MontFerret/ferret/v2/pkg/runtime"

// RegisterLib register `PATH` namespace functions.
// @namespace PATH
func RegisterLib(ns runtime.Namespace) {
	ns.Function().A1().
		Add("base", Base).
		Add("clean", Clean).
		Add("dir", Dir).
		Add("ext", Ext).
		Add("is_abs", IsAbs).
		Add("separate", Separate)

	ns.Function().A2().
		Add("match", Match)

	ns.Function().Var().
		Add("join", Join)
}
