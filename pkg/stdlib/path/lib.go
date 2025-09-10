package path

import "github.com/MontFerret/ferret/pkg/runtime"

// RegisterLib register `PATH` namespace functions.
// @namespace PATH
func RegisterLib(ns runtime.Namespace) error {
	ns.Functions().
		Set("BASE", Base).
		Set("CLEAN", Clean).
		Set("DIR", Dir).
		Set("EXT", Ext).
		Set("IS_ABS", IsAbs).
		Set("JOIN", Join).
		Set("MATCH", Match).
		Set("SEPARATE", Separate)

	return nil
}
