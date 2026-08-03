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

	ns.Function().A2().
		Add("WRITE", write2)

	ns.Function().A3().
		Add("WRITE", write3)
}
