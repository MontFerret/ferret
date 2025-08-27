package fs

import (
	"github.com/MontFerret/ferret/pkg/runtime"
)

// RegisterLib register `FS` namespace functions.
// @namespace FS
func RegisterLib(ns runtime.Namespace) error {
	return ns.
		Namespace("FS").
		RegisterFunctions(
			runtime.NewFunctionsFromMap(map[string]runtime.Function{
				"READ":  Read,
				"WRITE": Write,
			}))
}
