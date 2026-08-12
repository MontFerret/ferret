package fs

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// RegisterLib register `fs` namespace functions.
// @namespace fs
func RegisterLib(ns runtime.Namespace) {
	ns = ns.Namespace("fs")

	ns.Function().A1().
		Add("read", Read)

	ns.Function().A2().
		Add("write", write2)

	ns.Function().A3().
		Add("write", write3)
}
