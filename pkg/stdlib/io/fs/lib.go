package fs

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// RegisterLib register `FS` namespace functions.
// @namespace FS
func RegisterLib(ns runtime.Namespace) {
	ns = ns.Namespace("FS")

	ns.Function().A1().
		Add("READ", Read)

	ns.Function().Var().
		Add("WRITE", Write)
}
