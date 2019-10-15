package utils

import "github.com/MontFerret/ferret/pkg/runtime/core"

func RegisterLib(ns core.Namespace) error {
	return ns.RegisterFunctions(core.FunctionsMap{
		"WAIT":  Wait,
		"PRINT": Print,
	})
}
