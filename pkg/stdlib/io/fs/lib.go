package fs

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// RegisterLib register `FS` namespace functions.
// @namespace FS
func RegisterLib(ns core.Namespace) error {
	return ns.
		Namespace("FS").
		RegisterFunctions(
			core.NewFunctionsFromMap(map[string]core.Function{
				"READ":  Read,
				"WRITE": Write,
			}))
}
