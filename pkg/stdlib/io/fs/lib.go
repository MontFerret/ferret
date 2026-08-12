package fs

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// RegisterLib register `FS` namespace functions.
// @namespace FS
func RegisterLib(ns runtime.Namespace) {
	ns = ns.Namespace("FS")

	ns.Function().A1().
		Add("read", Read)

	ns.Function().A2().
		Add("write", write2)

	ns.Function().A3().
		Add("write", write3)
}
